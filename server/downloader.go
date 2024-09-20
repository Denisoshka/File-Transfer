package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"lab2"
	"lab2/requests"
	"lab2/utils"
	"net"
	"os"
	"path"
	"time"
)

const (
	BufSize = 1024 * 1024
)

type TCPDownloader struct {
	dirPath                string
	maxConnInactivityDelay time.Duration
	conn                   *net.TCPConn
	//buf      []byte
}

func NewTCPDownloader(dirPath string, maxDelay time.Duration,
	conn *net.TCPConn) *TCPDownloader {
	return &TCPDownloader{
		dirPath:                dirPath,
		maxConnInactivityDelay: maxDelay,
		conn:                   conn,
	}
}

func (u *TCPDownloader) Launch() {
	defer func(conn *net.TCPConn) { _ = conn.Close() }(u.conn)
	req, noticeConnection, err := u.handleInitialReq()
	if err != nil {
		if noticeConnection {
			_ = u.noticeDownloadFailed(err.Error())
		}
		lab2.Log.Info(
			"upload failed with error: ", err, " from ", u.conn.RemoteAddr(),
		)
	}

	filePath := path.Join(u.dirPath, req.Name)
	file, fileExists, err := utils.PrepareFile(filePath, req.DataSize)
	if err != nil {
		if fileExists {
			_ = u.noticeDownloadFailed("this file already exists")
		} else {
			_ = u.noticeDownloadFailed("unable to save file")
		}
		return
	}

	defer func(file *os.File) { _ = file.Close() }(file)

	total, err := u.downloadFile(req.DataSize, file)

	if total == req.DataSize {
		err = u.noticeDownloadSuccessful(fmt.Sprint("recorded ", total, " bytes"))
		lab2.Log.Debugln(u.conn, err)
		return
	}

	err = u.noticeDownloadFailed(fmt.Sprint("recorded ", total, " bytes"))
	lab2.Log.Debugln(u.conn, err)
}

func (u *TCPDownloader) downloadFile(dataSize int64, file *os.File) (total int64,
	err error) {
	bufManager := utils.NewBufferManager(BufSize)
	defer func(bufManager *utils.BufferManager) { bufManager.Close() }(bufManager)
	go fileWriter(file, bufManager)

	total = int64(0)
	for total < dataSize {
		buf, _ := bufManager.GetEmptyBuffer()
		var n int
		n, err = utils.ConnReadN(
			u.conn, buf.Data, buf.MaxCapacity, u.maxConnInactivityDelay,
		)
		buf.CurCapacity = n

		if err != nil {
			if err == io.EOF {
				break
			}

			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				err = u.noticeDownloadFailed(err.Error())
				lab2.Log.Debugln(u.conn, err)
				return total, err
			}
			lab2.Log.Errorln(u.conn, err)
			return total, err
		}

		bufManager.PushFullBuffer(buf)
		total += int64(n)
	}
	return total, nil
}

func fileWriter(file *os.File, bufManager *utils.BufferManager) {
	for {
		buf, opened := bufManager.GetFullBuffer()
		if !opened {
			return
		}
		_, err := utils.FileWriteN(file, buf.Data, buf.MaxCapacity)
		if err != nil {
			lab2.Log.Errorln("file: ", file.Name(), " error occurred ", err)
			return
		}
	}
}

func (u *TCPDownloader) handleInitialReq() (req *requests.Initial,
	noticeConnection bool, err error) {
	reqSizeBuf := make([]byte, 4)
	_, err = utils.ConnReadN(u.conn, reqSizeBuf, 4, u.maxConnInactivityDelay)
	if err != nil {
		return nil, false, err
	}

	initialReqSize := int32(binary.BigEndian.Uint32(reqSizeBuf))
	maxInitReqSize := requests.InitialSize(requests.MaxFileNameSize)
	if initialReqSize > maxInitReqSize {
		return nil, true, requests.IncorrectRequestSize
	}

	buf := make([]byte, initialReqSize)
	initialReq := &requests.Initial{}
	_, err = utils.ConnReadN(
		u.conn, buf, int(initialReqSize), u.maxConnInactivityDelay,
	)
	if err != nil {
		return nil, false, err
	}

	err = initialReq.DecodeFrom(buf)
	if err != nil {
		return nil, true, err
	}

	return initialReq, false, nil
}

func (u *TCPDownloader) noticeDownloadSuccessful(message string) (err error) {
	return notice(message, requests.SuccessResponse, u.conn)
}

func (u *TCPDownloader) noticeDownloadFailed(message string) (err error) {
	return notice(message, requests.ErrorResponse, u.conn)
}

func notice(message string, responseType int16, conn net.Conn) (err error) {
	if uint64(len(message)) > uint64(requests.MaxMessageSize) {
		message = ""
	}
	req, err := requests.NewResponse(responseType, message)
	if err != nil {
		lab2.Log.Debugln("what the fuck???")
		return err
	}
	data := make([]byte, req.HeaderSize)
	err = req.CodeTo(data)
	if err != nil {
		return err
	}

	_, err = utils.ConnWriteN(conn, data, int(req.HeaderSize))
	if err != nil {
		return err
	}
	return nil
}

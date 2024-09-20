package server

import (
	"encoding/binary"
	"errors"
	"fmt"
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

type Uploader struct {
	dirPath                string
	maxConnInactivityDelay time.Duration
	conn                   *net.TCPConn
	//buf      []byte
}

func NewUploader(dirPath string, maxDelay time.Duration,
	conn *net.TCPConn) *Uploader {
	return &Uploader{
		dirPath:                dirPath,
		maxConnInactivityDelay: maxDelay,
		conn:                   conn,
	}
}

func (u *Uploader) Launch() {
	defer func(conn *net.TCPConn) { _ = conn.Close() }(u.conn)
	req, noticeConnection, err := u.handleInitialReq()
	if err != nil {
		if noticeConnection {
			_ = u.noticeUploadFailed(err.Error())
		}
		lab2.Log.Info(
			"upload failed with error: ", err, " from ", u.conn.RemoteAddr(),
		)
	}
	filePath := path.Join(u.dirPath, req.Name)
	file, fileExists, err := utils.PrepareFile(filePath, int64(req.DataSize))
	if err != nil {
		if fileExists {
			_ = u.noticeUploadFailed("this file already exists")
		} else {
			_ = u.noticeUploadFailed("unable to save file")
		}
		return
	}

	bufManager := utils.NewBufferManager(BufSize)
	defer func(file *os.File) { _ = file.Close() }(file)
	defer func(bufManager *utils.BufferManager) { bufManager.Close() }(bufManager)
	go fileWriter(file, bufManager)

	total := int64(0)
	for total < req.DataSize {
		buf, _ := bufManager.GetEmptyBuffer()
		n, err := utils.ConnReadN(
			u.conn, buf.Data, buf.Total, u.maxConnInactivityDelay,
		)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				err = u.noticeUploadFailed(err.Error())
				lab2.Log.Debugln(u.conn, err)
			}
			lab2.Log.Errorln(u.conn, err)
			return
		}
		bufManager.PushFullBuffer(buf)
		total += int64(n)
	}
	if total == req.DataSize {
		err = u.noticeUploadSuccesfule(fmt.Sprint("recorded ", total, " bytes"))
		lab2.Log.Debugln(u.conn, err)
		return
	}
	err = u.noticeUploadFailed(fmt.Sprint("recorded ", total, " bytes"))
	lab2.Log.Debugln(u.conn, err)
}

func fileWriter(file *os.File, bufManager *utils.BufferManager) {
	for {
		buf, opened := bufManager.GetFullBuffer()
		if !opened {
			return
		}
		_, err := utils.FileWriteN(file, buf.Data, buf.Total)
		if err != nil {
			return
		}
	}
}

func (u *Uploader) handleInitialReq() (req *requests.Initial,
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

func (u *Uploader) noticeUploadSuccesfule(message string) (err error) {
	return notice(message, requests.SuccessResponse, u.conn)
}

func (u *Uploader) noticeUploadFailed(message string) (err error) {
	return notice(message, requests.ErrorResponse, u.conn)
}

func notice(message string, responseType int16, conn net.Conn) (err error) {
	req, err := requests.NewResponse(responseType, message)
	if err != nil {
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

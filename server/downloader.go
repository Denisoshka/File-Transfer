package server

import (
	"errors"
	"fmt"
	"lab2/utils"
	"lab2/utils/requests"
	"net"
	"os"
	"path"
	"time"
)

const (
	BufSize = 1024
)

type TCPDownloader struct {
	dirPath                string
	maxConnInactivityDelay time.Duration
	conn                   *net.TCPConn
	tracker                *SpeedTracker
}

func NewTCPDownloader(dirPath string, maxDelay time.Duration,
	conn *net.TCPConn, tracker *SpeedTracker) *TCPDownloader {
	return &TCPDownloader{
		dirPath:                dirPath,
		maxConnInactivityDelay: maxDelay,
		conn:                   conn,
		tracker:                tracker,
	}
}

func (u *TCPDownloader) Launch() (err error) {
	defer func(conn *net.TCPConn) { _ = conn.Close() }(u.conn)
	var req = &requests.Initial{}

	err = utils.ReadRequest(
		u.conn, u.maxConnInactivityDelay, req,
		int64(requests.InitialSize(requests.MaxFileNameSize)), nil,
	)
	//req, noticeConnection, err := u.handleInitialReq()
	LOG.Debugln(req)
	if err != nil {
		//if noticeConnection {
		_ = u.noticeDownloadFailed(err.Error())
		//}

		LOG.Info(
			"upload failed with error: ", err, " from ", u.conn.RemoteAddr(),
		)
		return
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

	total, err := u.fetchFile(req.DataSize, file)
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		//lab2.Log.Errorln(u.conn, err)
		err = u.noticeDownloadFailed(err.Error())
		return
	}
	if total == req.DataSize {
		err = u.noticeDownloadSuccessful(fmt.Sprint("recorded ", total, " bytes"))
		//lab2.Log.Debugln(u.conn, err)
		return
	}
	err = u.noticeDownloadFailed(fmt.Sprint("recorded ", total, " bytes"))
	//lab2.Log.Debugln(u.conn, err)
	return err
}

func (u *TCPDownloader) fetchFile(dataSize int64, file *os.File) (total int64, err error) {
	bufManager := utils.NewBufferManager(BufSize)
	tag := u.conn.RemoteAddr().String()
	//speedInfo := u.tracker.AddConnection(tag)
	//defer func(bufManager *utils.BufferManager) { bufManager.Close() }(bufManager)
	var funcErr error
	go func() {
		funcErr = fileWriter(file, bufManager)
	}()

	total = int64(0)
	for total < dataSize {
		buf, open := bufManager.GetForConsumer()
		if !open {
			break
		}

		var n int
		n, err = utils.ConnReadN(
			u.conn, buf.Data, buf.MaxCapacity, u.maxConnInactivityDelay,
		)
		LOG.Debugln("receive: ", n, " from ", tag)
		buf.CurCapacity = n
		if err != nil {
			break
		}

		bufManager.PushForConsumer(buf)
		total += int64(n)
	}
	bufManager.CloseForConsumer()

	if funcErr != nil {
		if err != nil {
			errors.Join(err, funcErr)
		} else {
			err = funcErr
		}
	}

	return total, err
}

func fileWriter(file *os.File, bufManager *utils.BufferManager) (err error) {
	defer func(bufManager *utils.BufferManager) { bufManager.CloseForPublisher() }(bufManager)
	for {
		buf, opened := bufManager.GetForPublisher()
		if !opened {
			return nil
		}
		_, err = utils.FileWriteN(file, buf.Data, buf.MaxCapacity)
		if err != nil {
			LOG.Errorln("file: ", file.Name(), " error occurred ", err)
			return err
		}
		bufManager.PushForPublisher(buf)
	}
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
		LOG.Debugln("what the fuck???")
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

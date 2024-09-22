package server

import (
	"errors"
	"fmt"
	serverinterfaces "lab2/server/server-interfaces"
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
	serverinterfaces.AbstractTCPDownloader
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
	if err != nil {
		LOG.Info(
			"upload failed with error: ", err, " from ", u.conn.RemoteAddr(),
		)
		err = u.noticeDownloadFailed(err.Error())
		return
	}

	filePath := path.Join(u.dirPath, req.Name)
	file, fileExists, err := prepareFile(filePath, req.DataSize)
	if err != nil {
		if fileExists {
			err = u.noticeDownloadFailed("this file already exists")
		} else {
			err = u.noticeDownloadFailed("unable to save file")
		}
		return
	}

	total, err := u.fetchFile(req.DataSize, file)
	_ = file.Close()

	if err != nil {
		LOG.Info("file fetch finished with", err)
		err = u.noticeDownloadFailed(err.Error())
		return
	}

	err = u.onDownloadEnd(total, req)
	return err
}

func (u *TCPDownloader) onDownloadEnd(total int64,
	req *requests.Initial) (err error) {
	if total == req.DataSize {
		msg := fmt.Sprint("recorded ", total, " bytes")
		err = u.noticeDownloadSuccessful(msg)
		if err != nil {
			LOG.Errorln(u.conn, err)
		} else {
			LOG.Info(u.conn, err)
		}
	} else {
		msg := fmt.Sprint("recorded ", total, " bytes")
		err = u.noticeDownloadFailed(msg)
	}
	return err
}

func (u *TCPDownloader) fetchFile(dataSize int64, file *os.File) (total int64, err error) {
	bufManager := utils.NewBufferManager(BufSize)
	tag := u.conn.RemoteAddr().String()
	connInfo := u.tracker.AddConnection(tag)

	var funcErr error
	go func() { funcErr = fileWriter(file, bufManager) }()

	total = int64(0)
	for received := 0; total < dataSize; {
		buf, open := bufManager.GetForConsumer()
		if !open {
			break
		}

		received, err = utils.ConnReadN(
			u.conn, buf.Data(),
			int(getReadQ(dataSize, total, int64(buf.MaxCapacity()))),
			u.maxConnInactivityDelay,
		)

		connInfo.AddRecordedQ(uint64(received))
		total += int64(received)

		if err != nil {
			break
		}

		bufManager.PushForConsumer(buf)
	}

	err = joinErrors(funcErr, err)
	connInfo.MarkAsExpired()
	bufManager.CloseForConsumer()
	LOG.Infoln("receive:", total, "from", tag, "with err", err)

	return total, err
}

func getReadQ(dataSize int64, total int64, expectedMaximum int64) int64 {
	remains := dataSize - total
	if remains > expectedMaximum {
		return expectedMaximum
	} else {
		return remains
	}
}

func fileWriter(file *os.File, bufManager *utils.BufferManager) (err error) {
	defer func(bufManager *utils.BufferManager) { bufManager.CloseForPublisher() }(bufManager)
	for {
		buf, opened := bufManager.GetForPublisher()
		if !opened {
			return nil
		}
		_, err = utils.FileWriteN(file, buf.Data(), buf.MaxCapacity())
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

func joinErrors(funcErr error, err error) error {
	if funcErr != nil {
		if err != nil {
			errors.Join(err, funcErr)
		} else {
			err = funcErr
		}
	}
	return err
}

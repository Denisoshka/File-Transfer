package client

import (
	"errors"
	"fmt"
	"io"
	"lab2/utils"
	"lab2/utils/requests"
	"net"
	"os"
	"time"
)

const (
	BufSize = 1024 * 4
)

type TCPDownloader struct {
	addr          *net.TCPAddr
	maxInactivity time.Duration
}

func NewTCPUploader(addr *net.TCPAddr, maxInactivity time.Duration) *TCPDownloader {
	return &TCPDownloader{
		addr:          addr,
		maxInactivity: maxInactivity,
	}
}

func (d *TCPDownloader) Launch(filePath string, uploadName string) (err error) {
	conn, err := net.DialTCP("tcp", nil, d.addr)
	if err != nil {
		return err
	}
	defer func(conn *net.TCPConn) { _ = conn.Close() }(conn)

	file, err := GetRegularFile(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	err = sendInitial(conn, stat, uploadName)
	if err != nil {
		return err
	}
	fmt.Println("send initial request for upload ", stat.Size(), " bytes")

	err = d.uploadFile(file, conn)
	err2 := d.handleResponse(conn)

	return errors.Join(err, err2)
}

func (d *TCPDownloader) handleResponse(conn *net.TCPConn) (err error) {
	req, uploadEndErr := d.onUploadEnd(conn)
	if uploadEndErr != nil {
		return errors.Join(err, uploadEndErr)
	}
	isSuccess := req.Status == requests.SuccessResponse
	msg := fmt.Sprint(
		"server response {status: ", isSuccess, " message: ", req.Message, "}",
	)
	fmt.Println(msg)
	return nil
}

func (d *TCPDownloader) uploadFile(file *os.File,
	conn *net.TCPConn) (err error) {
	bManager := utils.NewBufferManager(BufSize)
	var readerErr error
	go func() { readerErr = fileReader(bManager, file) }()
	total := int64(0)
	for q := 0; ; total += int64(q) {
		buf, opened := bManager.GetForPublisher()
		if !opened {
			err = readerErr
			break
		}

		q, err = utils.ConnWriteN(conn, buf.Data(), buf.CurCapacity())
		//fmt.Println("send ", q, " bytes")
		if err != nil {
			break
		}
		bManager.PushForPublisher(buf)
	}
	bManager.CloseForPublisher()

	if readerErr != nil {
		if err != nil {
			errors.Join(err, readerErr)
		} else {
			err = readerErr
		}
	}

	return err
}

func sendInitial(conn *net.TCPConn, stat os.FileInfo,
	uploadName string) (err error) {
	req, err := requests.NewInitial(stat.Size(), uploadName)
	if err != nil {
		return err
	}

	reqBuf := make([]byte, req.Size())
	err = req.CodeTo(reqBuf)
	if err != nil {
		return err
	}

	_, err = utils.ConnWriteN(conn, reqBuf, len(reqBuf))
	if err != nil {
		return err
	}
	return nil
}

func fileReader(bManager *utils.BufferManager, file *os.File) (err error) {
	for {
		buf, opened := bManager.GetForConsumer()
		if !opened {
			break
		}
		var n int
		n, err = utils.FileReadN(file, buf.Data(), buf.MaxCapacity())
		buf.SetCurCapacity(n)
		if err != nil {
			if err == io.EOF {
				bManager.PushForConsumer(buf)
				err = nil
			}
			break
		}
		bManager.PushForConsumer(buf)
	}
	bManager.CloseForConsumer()
	return err
}

func (d *TCPDownloader) readResponse(conn *net.TCPConn) (req *requests.Response, err error) {
	req = &requests.Response{}
	err = utils.ReadRequest(
		conn, d.maxInactivity, req,
		int64(requests.ResponseSize(requests.MaxMessageSize)), nil,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (d *TCPDownloader) onUploadEnd(conn *net.TCPConn) (req *requests.Response, err error) {
	req, err = d.readResponse(conn)
	if err != nil {
		return nil, err
	}
	return req, nil
}

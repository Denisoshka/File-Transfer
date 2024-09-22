package client

import (
	"encoding/binary"
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
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	err = sendInitial(conn, stat, uploadName)
	if err != nil {
		return err
	}

	bManager := utils.NewBufferManager(BufSize)
	var errPtr *error
	go func() { *errPtr = fileReader(bManager, file) }()
	for {
		buf, opened := bManager.GetForConsumer()
		if !opened {
			err = *errPtr
			break
		}
		_, err = utils.ConnWriteN(conn, buf.Data, buf.MaxCapacity)
		if err != nil {
			break
		}
		bManager.PushForConsumer(buf)
	}

	return d.onUploadEnd(conn)
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
	defer bManager.CloseConsumer()
	for {
		buf, opened := bManager.GetForConsumer()
		if !opened {
			return nil
		}
		var n int
		n, err = utils.FileReadN(file, buf.Data, buf.MaxCapacity)
		buf.CurCapacity = n
		if err != nil {
			if err == io.EOF {
				bManager.PushForConsumer(buf)
				err = nil
			}
			break
		}
		bManager.PushForConsumer(buf)
	}
	return err
}

func (d *TCPDownloader) readResponse(conn *net.TCPConn) (r *requests.Response, err error) {
	reqSizeBuf := make([]byte, 4)
	_, err = utils.ConnReadN(conn, reqSizeBuf, 4, d.maxInactivity)
	if err != nil {
		return nil, err
	}
	c := binary.BigEndian
	reqSize := c.Uint32(reqSizeBuf)
	buf := make([]byte, reqSize)
	c.PutUint32(buf[0:4], reqSize)
	_, err = utils.ConnReadN(conn, buf, int(reqSize), d.maxInactivity)
	if err != nil {
		return nil, err
	}
	r = &requests.Response{}
	err = r.DecodeFrom(buf)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (d *TCPDownloader) onUploadEnd(conn *net.TCPConn) (err error) {
	var req *requests.Response
	req, err = d.readResponse(conn)
	if err != nil {
		return err
	}
	isSuccess := req.ReqType == requests.SuccessResponse
	fmt.Println(
		"upload finished with success:", isSuccess, " and message ",
		req.Message,
	)
	return nil
}

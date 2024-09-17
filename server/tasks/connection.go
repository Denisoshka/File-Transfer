package tasks

import (
	"lab2/utils"
	"lab2/utils/requests"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	bufSize           = 1024
	maximumInactivity = time.Duration(time.Second * 10)
)

type Connection struct {
	dirPath string
	buf     []byte
	conn    *net.TCPConn
}

func NewConnection(dirPath string, conn *net.TCPConn) *Connection {
	return &Connection{
		dirPath: dirPath,
		conn:    conn,
	}
}

func (c *Connection) StartNetWorker() {
	conn := c.conn
	defer func(conn *net.TCPConn) {
		_ = conn.Close()
	}(conn)
	header, err := utils.ReadHeader(c)
	if err != nil {
		//todo
		return
	}

	filePath := filepath.Join(c.dirPath, header.FileName)
	file, err := utils.CreateFile(filePath, int64(header.FileSize))
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if err != nil {
		//todo
		return
	}

	buf1 := utils.NewBuffer(bufSize)
	buf2 := utils.NewBuffer(bufSize)
	getCh := make(chan *utils.Buffer, 2)
	putCh := make(chan *utils.Buffer, 2)
	errCh := make(chan error, 1)
	go startFileWorker(file, putCh, getCh, errCh)
	defer close(getCh)
	defer close(putCh)
	defer close(errCh)
	putCh <- buf1
	putCh <- buf2

	var totalRead uint64 = 0
	for n := 0; totalRead < header.FileSize; totalRead += uint64(n) {
		buf := <-getCh
		n, err = conn.Read(buf.Data)
		buf.Load = n
		if err != nil {
			break
		}
		putCh <- buf
	}
	close(putCh)
	close(errCh)
	err = <-errCh
	if err == nil {
		err = c.completeSuccessfully(filePath)
	}
	if err != nil {
		_ = c.completeUnsuccessfully(err)
	}
}

func startFileWorker(
	file *os.File,
	getCh <-chan *utils.Buffer,
	putCh chan<- *utils.Buffer,
	errCh chan<- error,
) {
	buf1 := <-getCh
	buf2 := <-getCh
	putCh <- buf1
	putCh <- buf2
	var err error = nil
	defer func(ch chan<- error, err error) {
		ch <- err
	}(errCh, err)

	for {
		buf, opened := <-getCh
		if !opened {
			break
		}
		err = utils.FileWriteBuffer(file, buf)
		if err != nil {
			break
		}
		putCh <- buf
	}
}

func (c *Connection) completeSuccessfully(filepath string) error {
	req, _ := requests.NewResponseReq(
		"complete successfully",
		requests.SuccessReq,
	)
	data, _ := req.MarshalBinary()
	c.conn.Write(data)
	return nil
}

func (c *Connection) completeUnsuccessfully(err error) error {
	req, err := requests.NewResponseReq(
		err.Error(),
		requests.ErrorReqMask,
	)
	if err != nil {
		req, _ = requests.NewResponseReq(
			"server error",
			requests.ErrorReqMask,
		)
	}
	data, _ := req.MarshalBinary()
	_, err = c.conn.Write(data)
	return nil
}

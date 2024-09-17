package tasks

import (
	"encoding/binary"
	"io"
	"net"
	"time"
	"xyi/helpers/requests"
)

var (
	bufSize           = 1024
	maximumInactivity = time.Duration(time.Second * 10)
)

type Connection struct {
	buf  []byte
	conn net.Conn
}

func NewConnection(conn *net.TCPConn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) StartNetWorker() {
	conn := c.conn
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	var headerSize uint32
	err := binary.Read(conn, binary.BigEndian, &headerSize)
	if err != nil {
		//todo
		return
	}
	headerBuf := make([]byte, headerSize)
	err = c.readN(headerBuf, int(headerSize))
	if err != nil {
		//todo
		return
	}
	header, err := requests.NewInitialRequest(headerBuf)
	if err != nil {
		//todo
		return
	}

	buf := make([]byte, bufSize)
	totalRead := uint64(0)
	for n := 0; totalRead < header.FileSize; totalRead += uint64(n) {
		n, err = conn.Read(buf)
		if err != nil {
			//todo
			return
		}
	}

}

func (c *Connection) readN(b []byte, n int) error {
	totalRead := 0
	conn := c.conn

	err := conn.SetDeadline(time.Now().Add(maximumInactivity))
	if err != nil {
		return err
	}

	for totalRead < n {
		bytesRead, err := conn.Read(b[totalRead:])
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		totalRead += bytesRead
	}

	return nil
}

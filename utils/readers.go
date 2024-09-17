package utils

import (
	"encoding/binary"
	"io"
	"lab2/server/tasks"
	"lab2/utils/requests"
	"net"
	"os"
	"time"
)

func CreateFile(fileName string, fileSize int64) (*os.File, error) {
	if _, err := os.Stat(fileName); os.IsExist(err) {
		return nil, err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	err = f.Truncate(0)
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	err = f.Truncate(fileSize)
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func ConnReadN(conn *net.TCPConn, b []byte, n int) error {
	totalRead := 0

	err := conn.SetDeadline(time.Now().Add(tasks.maximumInactivity))
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

func FileWriteBuffer(file *os.File, buffer *Buffer) error {
	for start := 0; start < buffer.Load; {
		ret, err := file.Write(buffer.Data[start:buffer.Load])
		if err != nil {
			return err
		}
		start += ret
	}
	return nil
}

func ReadHeader(c *tasks.Connection) (req *requests.InitialRequest, err error) {
	req = &requests.InitialRequest{}
	err = binary.Read(c.conn, binary.BigEndian, &req.HeaderSize)
	if err != nil {
		return
	}
	data := make([]byte, req.HeaderSize)
	binary.BigEndian.PutUint32(data[:4], req.HeaderSize)
	err = ConnReadN(c, data[4:], int(req.HeaderSize-4))
	if err != nil {
		return
	}
	err = req.UnmarshallBinaryTo(&data)
	if err != nil {
		return
	}
	return req, nil
}

/*
func ReadHeader(c *Connection) (req *requests.InitialRequest, err error) {
	var headerSize uint32
	err = binary.Read(c.conn, binary.BigEndian, &headerSize)
	if err != nil {
		//todo
		return nil, err
	}
	if headerSize < req.Size() {
		return nil, errors.New("header size is too small")
	}
	headerBuf := make([]byte, headerSize)
	binary.BigEndian.PutUint32(headerBuf[:4], headerSize)
	err = ConnReadN(c, headerBuf[4:], int(headerSize-4))
	if err != nil {
		return nil, err
	}
	err = req.UnmarshalBinary(headerBuf)
	if err != nil {
		return nil, err
	}
	return req, nil
}
*/

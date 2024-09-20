package client

import (
	"encoding/binary"
	"io"
	"lab2/requests"
	"lab2/utils"
	"net"
	"os"
	"time"
)

/*type DownloaderTask struct {
	fileUploadName string
	filePath       string
}*/

const (
	BufSize = 1024 * 4
)

type TCPDownloader struct {
	addr          *net.TCPAddr
	maxInactivity time.Duration
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
	fileReaderErrChan := make(chan error, 1)
	defer bManager.Close()

	go func() { _ = fileReader(bManager, file, fileReaderErrChan) }()
	for {
		buf, opened := bManager.GetFullBuffer()
		if !opened {
			break
		}
		_, err = utils.ConnWriteN(conn, buf.Data, buf.MaxCapacity)
		if err != nil {
			bManager.Close()
			break
		}
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

func fileReader(bManager *utils.BufferManager, file *os.File,
	errChan chan<- error) (err error) {
	defer bManager.Close()
	defer func(funcErr error) {
		if errChan != nil {
			errChan <- funcErr
		}
	}(err)

	for {
		buf, opened := bManager.GetEmptyBuffer()
		if !opened {
			return
		}
		var n int
		n, err = utils.FileReadN(file, buf.Data, buf.MaxCapacity)
		buf.CurCapacity = n
		if err != nil {
			if err == io.EOF {
				err = nil
				bManager.PushFullBuffer(buf)
			}
			break
		}
		bManager.PushFullBuffer(buf)
	}
	return
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
	_, err = d.readResponse(conn)
	if err != nil {
		return err
	}
	return nil
}

package server

import (
	"encoding/binary"
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
	dirPath  string
	maxDelay time.Duration
	conn     *net.TCPConn
	//buf      []byte
}

func NewUploader(dirPath string, maxDelay time.Duration,
	conn *net.TCPConn) *Uploader {
	return &Uploader{
		dirPath:  dirPath,
		maxDelay: maxDelay,
		conn:     conn,
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
	defer func(file *os.File) { _ = file.Close() }(file)
	buf := make()
	for total := int64(0); total < int64(req.DataSize); {

	}
}

func (u *Uploader) handleInitialReq() (req *requests.Initial,
	noticeConnection bool, err error) {
	reqSizeBuf := make([]byte, 4)
	_, err = utils.ConnReadN(u.conn, reqSizeBuf, 4)
	if err != nil {
		return nil, false, err
	}
	initialReqSize := binary.BigEndian.Uint32(reqSizeBuf)
	maxInitReqSize := requests.InitialSize(requests.MaxFileNameSize)
	if initialReqSize > maxInitReqSize {
		return nil, true, requests.IncorrectRequestSize
	}
	buf := make([]byte, initialReqSize)
	initialReq := &requests.Initial{}
	_, err = utils.ConnReadN(u.conn, buf, initialReqSize)
	if err != nil {
		return nil, false, err
	}
	err = initialReq.DecodeFrom(buf)
	if err != nil {
		return nil, true, err
	}
	return initialReq, false, nil
}

func (u *Uploader) noticeUploadFailed(message string) (err error) {
	req, err := requests.NewResponse(requests.ErrorResponse, message)
	if err != nil {
		return err
	}
	data := make([]byte, req.HeaderSize)
	_ = req.CodeTo(data)
	_, _ = utils.ConnWriteN(u.conn, data, req.HeaderSize)
	return nil
}

package requests

import (
	"errors"
	"net"
)

const InitialReq uint16 = 0
const MaxFileNameSize uint16 = 4096
const MaxFileSize int64 = 1024 * 1024 * 1024 * 1024
const ResponseReq uint16 = 1
const MaxMessageSize uint16 = 4096
const ErrorResponse = uint16(1)

const SuccessResponse = uint16(0)

var (
	MessageTooLargeError   = errors.New("message too long")
	BufferTooSmallError    = errors.New("buffer too small")
	InvalidHeaderSizeError = errors.New("invalid header size")
	FileNameTooLargeError  = errors.New("file name too large")
	FileSizeTooLargeError  = errors.New("file size too large")
	IncorrectRequestSize   = errors.New("incorrect request size")
)

func ReadRequest(conn net.Conn, n uint32, data []byte) {
}

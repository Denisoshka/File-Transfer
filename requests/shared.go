package requests

import (
	"errors"
)

const InitialReq int16 = 0
const MaxFileNameSize int16 = 4096
const MaxFileSize int64 = 1024 * 1024 * 1024 * 1024

const ResponseReq int16 = 1
const MaxMessageSize int16 = 4096
const ErrorResponse int16 = 1
const SuccessResponse int16 = 0

var (
	MessageTooLargeError   = errors.New("message too long")
	BufferTooSmallError    = errors.New("buffer too small")
	InvalidHeaderSizeError = errors.New("invalid header size")
	FileNameTooLargeError  = errors.New("file name too large")
	FileSizeTooLargeError  = errors.New("file size too large")
	IncorrectRequestSize   = errors.New("incorrect request size")
)

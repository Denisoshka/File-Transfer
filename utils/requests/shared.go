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
	MessageTooLargeError   = errors.New("long message")
	BufferTooSmallError    = errors.New("small buffer")
	InvalidHeaderSizeError = errors.New("invalid header size")
	FileNameTooLargeError  = errors.New("large file name")
	FileSizeTooLargeError  = errors.New("large file size")
	IncorrectRequestSize   = errors.New("incorrect request size")
)

package requests

import "errors"

var (
	IncorrectData         = errors.New("not enough data")
	FileNameTooLarge      = errors.New("file name too large")
	FileCheckSumWrongSize = errors.New("file checksum wrong HeaderSize")
)

const (
	MaxFileNameSize  = uint64(4048)
	MaxMessageSize   = uint64(4048)
	FileCheckSumSize = uint64(64)
	ErrorReqMask     = uint16(1)
	SuccessReq       = uint16(0)
)

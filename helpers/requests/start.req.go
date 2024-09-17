package requests

import (
	"encoding/binary"
	"errors"
)

var (
	InitialRequestSize = 8

	//requestSizeStart = 0
	//requestSizeEnd   = 4
	//
	//fileNameSizeStart = 4
	//fileNameSizeEnd   = 6
	//fileNameStart     = 6
	maxFileNameSize = 4096
)

type InitialRequest struct {
	HeaderSize   uint32
	FileSize     uint64
	FileNameSize uint16
	FileName     string
	FileChecksum []byte
}

func NewInitialRequest(data []byte) (*InitialRequest, error) {
	req := &InitialRequest{}
	return req, unmarshallLogic(req, data)
}

func (req *InitialRequest) UnmarshalBinary(data []byte) error {
	return unmarshallLogic(req, data)
}

func unmarshallLogic(req *InitialRequest, data []byte) error {
	encoder := binary.BigEndian
	if len(data) < InitialRequestSize {
		return errors.New("not enough data")
	}
	req.HeaderSize = encoder.Uint32(data[0:4])
	if req.HeaderSize > uint32(len(data)) {
		return errors.New("not enough data to unmarshall header")
	}

	req.FileSize = encoder.Uint64(data[4:12])
	req.FileNameSize = encoder.Uint16(data[12:14])
	if uint64(18+req.FileNameSize) > uint64(req.HeaderSize) {
		return errors.New("not enough data to unmarshall header")
	}
	req.FileName = string(data[14 : 14+req.FileNameSize])
	req.FileChecksum = data[14+req.FileNameSize : 22+req.FileNameSize]
	return nil
}

func (req *InitialRequest) MarshalBinary() ([]byte, error) {

}

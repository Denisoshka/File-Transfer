package requests

import (
	"encoding/binary"
	"errors"
)

type ResponseReq struct {
	HeaderSize  uint32
	Status      uint16
	MessageSize uint16
	Message     string
}

func (req *ResponseReq) Size() uint32 {
	return responseReqSize(req)
}

func NewResponseReq(message string, status uint16) (req *ResponseReq, err error) {
	messageSize := uint64(len(message))
	if messageSize > MaxMessageSize {
		return nil, errors.New("message messageSize too big")
	}
	req = &ResponseReq{
		HeaderSize:  uint32(messageSize + 8),
		Status:      status,
		MessageSize: uint16(messageSize),
		Message:     message,
	}
	return req, nil
}

func responseReqSize(req *ResponseReq) uint32 {
	return uint32(len(req.Message) + 8)
}

func (req *ResponseReq) MarshalBinary() ([]byte, error) {
	size := responseReqSize(req)
	data := make([]byte, size)
	return data, req.MarshallBinaryTo(&data)
}

func (req *ResponseReq) MarshallBinaryTo(data *[]byte) error {
	reqSize := responseReqSize(req)
	if reqSize != uint32(len(*data)) {
		return IncorrectData
	}
	decoder := binary.BigEndian
	decoder.PutUint32((*data)[0:4], reqSize)
	decoder.PutUint16((*data)[4:6], req.Status)
	decoder.PutUint16((*data)[6:8], uint16(len(req.Message)))
	copy((*data)[8:], req.Message)

	return nil
}

func (req *ResponseReq) UnmarshalBinary(data []byte) error {
	return req.UnmarshallBinaryTo(&data)
}

func (req *ResponseReq) UnmarshallBinaryTo(data *[]byte) error {
	if len(*data) < 8 {
		return IncorrectData
	}
	decoder := binary.BigEndian
	req.HeaderSize = decoder.Uint32((*data)[0:4])
	if req.HeaderSize != uint32(len(*data)) {
		return IncorrectData
	}
	req.Status = decoder.Uint16((*data)[4:6])
	req.MessageSize = decoder.Uint16((*data)[6:8])
	req.Message = string((*data)[8:])
	return nil
}

func IsErrorResponse(req *ResponseReq) bool {
	return (req.Status & ErrorReqMask) != 0
}

func (req *ResponseReq) IsErrorResponse() bool {
	return IsErrorResponse(req)
}

func (req *ResponseReq) IsSuccessResponse() bool {
	return IsSuccessResponse(req)
}
func IsSuccessResponse(req *ResponseReq) bool {
	return (req.Status & 1) == 0
}

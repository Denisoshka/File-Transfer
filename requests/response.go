package requests

import (
	"encoding/binary"
)

type Response struct {
	HeaderSize  uint32
	ReqType     uint16
	Status      uint16
	MessageSize uint16
	Message     string
}

// NewResponse
// err MessageTooLargeError if len(message) > MaxMessageSize
func NewResponse(status uint16, message string) (response *Response,
	err error) {
	messageSize := len(message)
	if uint64(messageSize) > uint64(MaxMessageSize) {
		return nil, MessageTooLargeError
	}

	return &Response{
		HeaderSize:  ResponseSize(uint16(messageSize)),
		ReqType:     ResponseReq,
		Status:      status,
		MessageSize: uint16(messageSize),
		Message:     message,
	}, nil
}

func ResponseSize(messageSize uint16) uint32 {
	return uint32(4 + 2 + 2 + 2 + messageSize)
}

// CodeTo
// err BufferTooSmallError if HeaderSize > len(data)
func (r *Response) CodeTo(data []byte) (err error) {
	size := r.HeaderSize
	n := len(data)
	if uint64(size) > uint64(n) {
		return BufferTooSmallError
	}
	decoder := binary.BigEndian
	decoder.PutUint32(data[0:4], r.HeaderSize)
	decoder.PutUint16(data[4:6], r.ReqType)
	decoder.PutUint16(data[6:8], r.Status)
	decoder.PutUint16(data[8:10], r.MessageSize)
	copy(data[10:r.MessageSize], r.Message)
	return nil
}

// DecodeFrom
// err BufferTooSmallError if HeaderSize > len(data)
// err InvalidHeaderSizeError if header have incorrect size
func (r *Response) DecodeFrom(data []byte) (err error) {
	n := len(data)
	if uint64(n) < uint64(ResponseSize(0)) {
		return BufferTooSmallError
	}
	decoder := binary.BigEndian
	r.HeaderSize = decoder.Uint32(data[0:4])
	r.ReqType = decoder.Uint16(data[4:6])
	r.Status = decoder.Uint16(data[6:8])
	r.MessageSize = decoder.Uint16(data[8:10])
	size := ResponseSize(r.MessageSize)
	if uint64(size) > uint64(n) {
		return BufferTooSmallError
	}
	if uint64(r.HeaderSize) != uint64(size) {
		return InvalidHeaderSizeError
	}
	r.Message = string(data[10:size])
	return nil
}

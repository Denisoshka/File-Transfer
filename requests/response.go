package requests

import (
	"encoding/binary"
)

type Response struct {
	HeaderSize  int32
	ReqType     int16
	Status      int16
	MessageSize int16
	Message     string
}

// NewResponse
// err MessageTooLargeError if len(message) > MaxMessageSize
func NewResponse(status int16, message string) (response *Response,
	err error) {
	messageSize := len(message)
	if int64(messageSize) > int64(MaxMessageSize) {
		return nil, MessageTooLargeError
	}

	return &Response{
		HeaderSize:  ResponseSize(int16(messageSize)),
		ReqType:     ResponseReq,
		Status:      status,
		MessageSize: int16(messageSize),
		Message:     message,
	}, nil
}

func ResponseSize(messageSize int16) int32 {
	return int32(4 + 2 + 2 + 2 + messageSize)
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
	decoder.PutUint32(data[0:4], uint32(r.HeaderSize))
	decoder.PutUint16(data[4:6], uint16(r.ReqType))
	decoder.PutUint16(data[6:8], uint16(r.Status))
	decoder.PutUint16(data[8:10], uint16(r.MessageSize))
	copy(data[10:r.MessageSize], r.Message)
	return nil
}

// DecodeFrom
// err BufferTooSmallError if HeaderSize > len(data)
// err InvalidHeaderSizeError if header have incorrect size
func (r *Response) DecodeFrom(data []byte) (err error) {
	n := len(data)
	if int64(n) < int64(ResponseSize(0)) {
		return BufferTooSmallError
	}
	decoder := binary.BigEndian
	r.HeaderSize = int32(decoder.Uint32(data[0:4]))
	r.ReqType = int16(decoder.Uint16(data[4:6]))
	r.Status = int16(decoder.Uint16(data[6:8]))
	r.MessageSize = int16(decoder.Uint16(data[8:10]))
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

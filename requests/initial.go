package requests

import "encoding/binary"

type Initial struct {
	HeaderSize int32
	ReqType    int16
	DataSize   int64
	NameSize   int16
	Name       string
}

// NewInitial
// err FileNameTooLargeError if len(name) > MaxFileNameSize
// err FileSizeTooLargeError if len(name) > MaxFileSize
func NewInitial(dataSize int64, name string) (req *Initial, err error) {
	nameSize := len(name)
	if int64(nameSize) > int64(MaxFileNameSize) {
		return nil, FileNameTooLargeError
	}
	if dataSize > MaxFileSize {
		return nil, FileSizeTooLargeError
	}

	return &Initial{
		HeaderSize: InitialSize(int16(nameSize)),
		ReqType:    InitialReq,
		DataSize:   dataSize,
		NameSize:   int16(nameSize),
		Name:       name,
	}, nil
}

func InitialSize(nameSize int16) int32 {
	return int32(4 + 2 + 8 + 2 + 2 + nameSize)
}

func (r *Initial) CodeTo(data []byte) (err error) {
	size := r.HeaderSize
	n := len(data)
	if uint64(size) > uint64(n) {
		return BufferTooSmallError
	}
	dc := binary.BigEndian
	dc.PutUint32(data[0:4], uint32(r.HeaderSize))
	dc.PutUint16(data[4:6], uint16(r.ReqType))
	dc.PutUint64(data[6:14], uint64(r.DataSize))
	dc.PutUint16(data[14:16], uint16(r.NameSize))
	copy(data[16:size], r.Name)
	return nil
}

func (r *Initial) DecodeFrom(data []byte) (err error) {
	n := len(data)
	if uint64(n) < uint64(InitialSize(0)) {
		return BufferTooSmallError
	}
	dc := binary.BigEndian
	r.HeaderSize = int32(dc.Uint32(data[0:4]))
	r.ReqType = int16(dc.Uint16(data[4:6]))
	r.DataSize = int64(dc.Uint64(data[6:14]))
	r.NameSize = int16(dc.Uint16(data[14:16]))
	size := InitialSize(r.NameSize)
	if uint64(size) > uint64(n) {
		return BufferTooSmallError
	}
	if uint64(r.HeaderSize) != uint64(size) {
		return InvalidHeaderSizeError
	}
	r.Name = string(data[16:size])
	return nil
}

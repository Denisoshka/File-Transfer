package requests

import "encoding/binary"

type Initial struct {
	HeaderSize uint32
	ReqType    uint16
	DataSize   uint64
	NameSize   uint16
	Name       string
}

// NewInitial
// err FileNameTooLargeError if len(name) > MaxFileNameSize
// err FileSizeTooLargeError if len(name) > MaxFileSize
func NewInitial(dataSize uint64, name string) (req *Initial, err error) {
	nameSize := len(name)
	if uint64(nameSize) > uint64(MaxFileNameSize) {
		return nil, FileNameTooLargeError
	}
	if dataSize > uint64(MaxFileSize) {
		return nil, FileSizeTooLargeError
	}

	return &Initial{
		HeaderSize: InitialSize(uint16(nameSize)),
		ReqType:    InitialReq,
		DataSize:   dataSize,
		NameSize:   uint16(nameSize),
		Name:       name,
	}, nil
}

func InitialSize(nameSize uint16) uint32 {
	return uint32(4 + 2 + 8 + 2 + 2 + nameSize)
}

func (r *Initial) CodeTo(data []byte) (err error) {
	size := r.HeaderSize
	n := len(data)
	if uint64(size) > uint64(n) {
		return BufferTooSmallError
	}
	dc := binary.BigEndian
	dc.PutUint32(data[0:4], r.HeaderSize)
	dc.PutUint16(data[4:6], r.ReqType)
	dc.PutUint64(data[6:14], r.DataSize)
	dc.PutUint16(data[14:16], r.NameSize)
	copy(data[16:size], r.Name)
	return nil
}

func (r *Initial) DecodeFrom(data []byte) (err error) {
	n := len(data)
	if uint64(n) < uint64(InitialSize(0)) {
		return BufferTooSmallError
	}
	dc := binary.BigEndian
	r.HeaderSize = dc.Uint32(data[0:4])
	r.ReqType = dc.Uint16(data[4:6])
	r.DataSize = dc.Uint64(data[6:14])
	r.NameSize = dc.Uint16(data[14:16])
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

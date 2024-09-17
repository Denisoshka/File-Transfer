package requests

import "encoding/binary"

type InitialRequest struct {
	HeaderSize   uint32
	FileSize     uint64
	FileNameSize uint16
	FileName     string
	FileCheckSum string
}

func (req *InitialRequest) Size() uint32 {
	return initialRequestSize(req)
}

func initialRequestSize(req *InitialRequest) uint32 {
	return uint32(4 + 8 + 2 + len(req.FileName) + 64)
}

func (req *InitialRequest) MarshalBinary() ([]byte, error) {
	data := make([]byte, initialRequestSize(req))
	return data, req.MarshallBinaryTo(&data)
}

func (req *InitialRequest) MarshallBinaryTo(data *[]byte) error {
	reqSize := initialRequestSize(req)
	if reqSize != uint32(len(*data)) {
		return IncorrectData
	}
	fileNameSize := uint32(len(req.FileName))
	if fileNameSize > MaxFileNameSize {
		return FileNameTooLarge
	}
	if uint32(len(req.FileCheckSum)) != FileCheckSumSize {
		return FileCheckSumWrongSize
	}

	decoder := binary.BigEndian
	decoder.PutUint32((*data)[0:4], reqSize)
	decoder.PutUint64((*data)[4:12], req.FileSize)
	decoder.PutUint16((*data)[12:14], uint16(fileNameSize))
	copy((*data)[14:14+fileNameSize], req.FileName)
	copy(
		(*data)[14+fileNameSize:14+fileNameSize+FileCheckSumSize], req.FileCheckSum,
	)
	return nil
}

func (req *InitialRequest) UnmarshalBinary(data []byte) error {
	return req.UnmarshallBinaryTo(&data)
}

func (req *InitialRequest) UnmarshallBinaryTo(data *[]byte) error {
	decoder := binary.BigEndian
	req.HeaderSize = decoder.Uint32((*data)[0:4])
	if req.HeaderSize != uint32(len(*data)) {
		return IncorrectData
	}
	req.FileSize = decoder.Uint64((*data)[4:12])
	req.FileNameSize = decoder.Uint16((*data)[12:14])
	req.FileName = string((*data)[14 : 14+req.FileNameSize])
	req.FileCheckSum = string(
		(*data)[14+req.FileNameSize : 14+uint32(req.FileNameSize)+FileCheckSumSize],
	)
	return nil
}

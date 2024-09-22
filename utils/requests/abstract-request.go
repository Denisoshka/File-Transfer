package requests

type AbstractRequest interface {
	Size() int32
	CodeTo(data []byte) (err error)
	DecodeFrom(data []byte) (err error)
}

package utils

type Buffer struct {
	Data []byte
	Load int
}

func (b *Buffer) SetLoad(load int) {
	b.Load = load
}

func NewBuffer(n int) *Buffer {
	return &Buffer{
		Data: make([]byte, n),
		Load: 0,
	}
}

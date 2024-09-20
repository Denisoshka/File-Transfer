package utils

const bufQ = 2

type BufferManager struct {
	forFile   chan *Buffer
	forSocket chan *Buffer
}

type Buffer struct {
	Data  []byte
	Total int
}

func NewBuffer(bufSize int) (b *Buffer) {
	return &Buffer{
		Data:  make([]byte, bufSize),
		Total: 0,
	}
}

func NewBufferManager(bufSize int) (m *BufferManager) {
	m = &BufferManager{
		forFile:   make(chan *Buffer, bufQ),
		forSocket: make(chan *Buffer, bufQ),
	}
	m.PushEmptyBuffer(NewBuffer(bufSize))
	m.PushEmptyBuffer(NewBuffer(bufSize))
	return m
}

func (m *BufferManager) Close() {
	close(m.forFile)
	close(m.forSocket)
}

func (m *BufferManager) GetEmptyBuffer() (b *Buffer, opened bool) {
	b, opened = <-m.forSocket
	return
}

func (m *BufferManager) PushFullBuffer(b *Buffer) {
	m.forFile <- b
}

func (m *BufferManager) GetFullBuffer() (b *Buffer, opened bool) {
	b, opened = <-m.forFile
	return
}

func (m *BufferManager) PushEmptyBuffer(b *Buffer) {
	m.forSocket <- b
}

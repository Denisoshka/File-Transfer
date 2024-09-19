package utils

type Buffer struct {
	Data []byte
	load int
}

const bufQ = 2

type BufferManager struct {
	//inited atomic.Bool
	forFile   chan *Buffer
	forSocket chan *Buffer
}

func NewBufferManager() (m *BufferManager) {
	return &BufferManager{
		forFile:   make(chan *Buffer, bufQ),
		forSocket: make(chan *Buffer, bufQ),
	}
}

func (m *BufferManager) close() {
	close(m.forFile)
	close(m.forSocket)
}

func (m *BufferManager) GetEmptyBuffer() *Buffer {
	return <-m.forSocket
}
func (m *BufferManager) PushFullBuffer(b *Buffer) {
	m.forFile <- b
}

func (m *BufferManager) GetFullBuffer() *Buffer {
	return <-m.forFile
}
func (m *BufferManager) PushEmptyBuffer(b *Buffer) {
	m.forSocket <- b
}

func (b *Buffer) SetLoad(load int) {
	b.load = load
}

func NewBuffer(n int) *Buffer {
	return &Buffer{
		Data: make([]byte, n),
		load: 0,
	}
}

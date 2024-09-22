package utils

const bufQ = 2

type BufferManager struct {
	consumer        chan *Buffer
	recipient       chan *Buffer
	forFileClosed   int8
	forSocketClosed int8
}

type Buffer struct {
	Data        []byte
	MaxCapacity int
	CurCapacity int
}

func NewBuffer(bufSize int) (b *Buffer) {
	return &Buffer{
		Data:        make([]byte, bufSize),
		MaxCapacity: bufSize,
		CurCapacity: 0,
	}
}

func NewBufferManager(bufSize int) (m *BufferManager) {
	m = &BufferManager{
		consumer:  make(chan *Buffer, bufQ),
		recipient: make(chan *Buffer, bufQ),
	}
	m.PushForPublisher(NewBuffer(bufSize))
	m.PushForPublisher(NewBuffer(bufSize))
	return m
}

func (m *BufferManager) CloseConsumer() {
	close(m.consumer)
}

func (m *BufferManager) ClosePublisher() {
	close(m.recipient)
}

func (m *BufferManager) GetForPublisher() (b *Buffer, opened bool) {
	b, opened = <-m.recipient
	return
}

func (m *BufferManager) PushForPublisher(b *Buffer) {
	m.consumer <- b
}

func (m *BufferManager) GetForConsumer() (b *Buffer, opened bool) {
	b, opened = <-m.consumer
	return
}

func (m *BufferManager) PushForConsumer(b *Buffer) {
	m.recipient <- b
}

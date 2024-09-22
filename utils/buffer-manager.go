package utils

const bufQ = 2

type BufferManager struct {
	consumer        chan *Buffer
	publisher       chan *Buffer
	forFileClosed   int8
	forSocketClosed int8
}

type Buffer struct {
	data        []byte
	maxCapacity int
	curCapacity int
}

func (b *Buffer) CurCapacity() int {
	return b.curCapacity
}

func (b *Buffer) SetCurCapacity(curCapacity int) {
	b.curCapacity = curCapacity
}

func (b *Buffer) Data() []byte {
	return b.data
}

func (b *Buffer) MaxCapacity() int {
	return b.maxCapacity
}

func NewBuffer(bufSize int) (b *Buffer) {
	return &Buffer{
		data:        make([]byte, bufSize),
		maxCapacity: bufSize,
		curCapacity: 0,
	}
}

func NewBufferManager(bufSize int) (m *BufferManager) {
	m = &BufferManager{
		consumer:  make(chan *Buffer, bufQ),
		publisher: make(chan *Buffer, bufQ),
	}
	m.PushForPublisher(NewBuffer(bufSize))
	m.PushForPublisher(NewBuffer(bufSize))
	return m
}

func (m *BufferManager) CloseForConsumer() {
	close(m.publisher)
}

func (m *BufferManager) CloseForPublisher() {
	close(m.consumer)
}

func (m *BufferManager) GetForPublisher() (b *Buffer, opened bool) {
	b, opened = <-m.publisher
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
	m.publisher <- b
}

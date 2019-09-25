package udwBytes

import (
	"sync"
)

type LimitedBuffer struct {
	maxSizeInBytes int
	buffer         []byte
	count          int
	lock           sync.Mutex
}

func NewLimitedBuffer(maxSizeInBytes int) *LimitedBuffer {
	if maxSizeInBytes <= 0 {
		panic("maxSizeInBytes must > 0")
	}
	buffer := &LimitedBuffer{
		maxSizeInBytes: maxSizeInBytes,
		buffer:         make([]byte, 0, maxSizeInBytes),
	}
	return buffer
}

func (buffer *LimitedBuffer) Reset() {
	buffer.lock.Lock()
	buffer.count = 0
	buffer.buffer = buffer.buffer[:0]
	buffer.lock.Unlock()
}

func (buffer *LimitedBuffer) AddBytes(bytesData []byte) {
	buffer.lock.Lock()
	defer buffer.lock.Unlock()
	dataCount := len(bytesData)
	if dataCount > buffer.maxSizeInBytes {
		buffer.buffer = bytesData[dataCount-buffer.maxSizeInBytes:]
		buffer.count = buffer.maxSizeInBytes
		return
	}
	delta := dataCount + buffer.count - buffer.maxSizeInBytes
	if delta > 0 {
		buffer.buffer = buffer.buffer[delta:]
		buffer.count -= delta
	}
	for _, b := range bytesData {
		buffer.buffer = append(buffer.buffer, b)
	}
	buffer.count += dataCount
}

func (buffer *LimitedBuffer) Add(data string) {
	buffer.AddBytes([]byte(data))
}

func (buffer *LimitedBuffer) AddLine(data string) {
	buffer.Add(data + "\n")
}

func (buffer *LimitedBuffer) ToString() string {
	buffer.lock.Lock()
	defer buffer.lock.Unlock()
	return string(buffer.buffer[:buffer.count])
}

package udwBytes

import (
	"bytes"
)

type QueueByteSlice struct {
	lenQueue QueueInt

	dataBuffer bytes.Buffer
}

func (rb *QueueByteSlice) Init() {
	dataBuffer := bytes.NewBuffer(make([]byte, 0, 100*1024))
	rb.dataBuffer = *dataBuffer
	rb.lenQueue.Init()

}

func (rb *QueueByteSlice) AddOne(b []byte) {
	rb.lenQueue.AddOne(len(b))

	rb.dataBuffer.Write(b)
}

func (rb *QueueByteSlice) GetOne() []byte {
	l := rb.lenQueue.GetOne()
	return rb.dataBuffer.Bytes()[:l]
}

func (rb *QueueByteSlice) RemoveOne() {
	l := rb.lenQueue.GetOne()
	rb.dataBuffer.Next(l)
	rb.lenQueue.RemoveOne()
}
func (rb *QueueByteSlice) HasData() bool {
	return rb.lenQueue.HasData()
}

package udwBytes

type QueueInt struct {
	buf []int
	off int
}

func (buf *QueueInt) Init() {
	buf.buf = make([]int, 0, 1024)
	buf.off = 0
}

func (q *QueueInt) AddOne(i int) {
	if q.off+1 > len(q.buf) {
		if q.off > len(q.buf)/2 {
			thisSize := len(q.buf) - q.off
			copy(q.buf, q.buf[q.off:])
			q.off = 0
			q.buf = q.buf[:thisSize]
		} else {
			newCap := len(q.buf) * 2
			if newCap < 32 {
				newCap = 32
			}
			thisSize := len(q.buf) - q.off
			newBuf := make([]int, newCap)
			copy(newBuf, q.buf[q.off:])
			q.off = 0
			q.buf = q.buf[:thisSize]
		}
	}
	q.buf = append(q.buf, i)
}
func (buf *QueueInt) GetOne() int {
	return buf.buf[buf.off]
}

func (buf *QueueInt) RemoveOne() {
	buf.off++
}
func (buf *QueueInt) HasData() bool {
	return len(buf.buf) > buf.off
}

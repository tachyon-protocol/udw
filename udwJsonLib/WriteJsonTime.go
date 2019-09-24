package udwJsonLib

import (
	"strconv"
	"time"
)

func WriteJsonTime(ctx *Context, t time.Time) {
	y := t.Year()
	if y < 0 || y >= 10000 {

		panic("WriteJsonTime: year outside of range [0,9999] " + strconv.Itoa(y))
	}
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	writerTryRealloc(ctx, len(time.RFC3339Nano)+2)
	WriterWriteByte(ctx, '"')
	writeIntWithWidth(ctx, year, 4)
	WriterWriteByte(ctx, '-')
	writeIntWithWidth(ctx, int(month), 2)
	WriterWriteByte(ctx, '-')
	writeIntWithWidth(ctx, day, 2)
	WriterWriteByte(ctx, 'T')
	writeIntWithWidth(ctx, hour, 2)
	WriterWriteByte(ctx, ':')
	writeIntWithWidth(ctx, minute, 2)
	WriterWriteByte(ctx, ':')
	writeIntWithWidth(ctx, second, 2)
	nonsecond := t.Nanosecond()
	if nonsecond > 0 {
		WriterWriteByte(ctx, '.')
		writeIntWithWidth(ctx, nonsecond, 9)
	}

	_, offset := t.Zone()
	if offset == 0 {
		WriterWriteByte(ctx, 'Z')
	} else {
		if offset < 0 {
			WriterWriteByte(ctx, '-')
			offset = -offset
		} else {
			WriterWriteByte(ctx, '+')
		}
		writeIntWithWidth(ctx, offset/3600, 2)
		WriterWriteByte(ctx, ':')
		writeIntWithWidth(ctx, (offset/60)%60, 2)
	}

	WriterWriteByte(ctx, '"')
	return
}

func writeIntWithWidth(ctx *Context, x int, width int) {
	u := uint(x)
	if x < 0 {
		WriterWriteByte(ctx, '-')
		u = uint(-x)
	}

	var buf [20]byte
	i := len(buf)
	for u >= 10 {
		i--
		q := u / 10
		buf[i] = byte('0' + u - q*10)
		u = q
	}
	i--
	buf[i] = byte('0' + u)

	for w := len(buf) - i; w < width; w++ {
		WriterWriteByte(ctx, '0')
	}
	WriterWriteByteList(ctx, buf[i:])
}

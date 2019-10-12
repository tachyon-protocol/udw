package udwGoReader

import (
	"sort"
	"strconv"
)

type FilePos struct {
	filePath string
	lines    []int
}

func NewPosFile(filepath string, content []byte) *FilePos {
	lines := make([]int, 0)
	r := NewReader(content, nil)
	lines = append(lines, 0)
	for {
		r.ReadUntilByte('\n')
		if r.IsEof() {
			break
		}
		lines = append(lines, r.pos)
	}
	return &FilePos{
		filePath: filepath,
		lines:    lines,
	}
}

func (p *FilePos) PosString(pos int) string {
	if p == nil {
		return "<nil>"
	}
	line := p.GetLineWithPos(pos)
	return p.filePath + ":" + strconv.Itoa(line)
}

func (p *FilePos) GetLineWithPos(pos int) int {
	return sort.SearchInts(p.lines, pos)
}

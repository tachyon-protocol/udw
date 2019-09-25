package udwStrings

type LineSpliter struct {
	content string
	posList []int
}

func NewLineSpliter(content string) *LineSpliter {
	return &LineSpliter{
		content: content,
		posList: getLinePosList(content),
	}
}

func (ls *LineSpliter) GetLineNumByPos(pos int) int {
	if pos <= 0 {
		return 1
	}
	for i := range ls.posList {
		if pos < ls.posList[i] {
			return i + 1
		}
	}
	return len(ls.posList) + 1
}

func (ls *LineSpliter) GetLineContent(lineNumber int) string {
	if lineNumber <= 0 {
		return ""
	}
	if lineNumber > len(ls.posList)+1 {
		return ""
	}
	if len(ls.posList) == 0 {
		return ls.content
	}
	if lineNumber == len(ls.posList)+1 {
		return ls.content[ls.posList[len(ls.posList)-1]:]
	}
	if lineNumber == 1 {
		return ls.content[0 : ls.posList[0]-1]
	}
	return ls.content[ls.posList[lineNumber-2] : ls.posList[lineNumber-1]-1]
}

func (ls *LineSpliter) GetTotalLineNumber() int {
	return len(ls.posList) + 1
}

func getLinePosList(content string) (out []int) {
	pos := 0
	l := len(content)
	for {
		if pos >= l {
			return
		}
		r := content[pos]
		if r == '\n' {
			out = append(out, pos+1)
		}
		pos++
	}
}

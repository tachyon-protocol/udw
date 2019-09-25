package udwFile

import "github.com/tachyon-protocol/udw/udwBytes"

func PathClean(path string) string {
	if len(path) == 0 {
		return path
	}
	outBuf := &udwBytes.BufWriter{}
	thisPartSize := 0
	isThisPartEmptyOrAllDot := true
	for i := 0; i < len(path); i++ {
		b := path[i]
		if b == '/' {
			if isThisPartEmptyOrAllDot {
				outBuf.AddPos(-thisPartSize)
			} else {
				outBuf.WriteByte(b)
			}
			thisPartSize = 0
			isThisPartEmptyOrAllDot = true
			continue
		}
		if b != '.' {
			isThisPartEmptyOrAllDot = false
		}
		outBuf.WriteByte(b)
		thisPartSize++
	}
	if isThisPartEmptyOrAllDot {
		outBuf.AddPos(-thisPartSize)
	}

	buf := outBuf.GetBytes()
	if len(buf) > 0 && buf[len(buf)-1] == '/' {
		return string(buf[:len(buf)-1])
	}
	return outBuf.GetString()
}

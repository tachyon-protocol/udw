package udwGoWriter

import "bytes"

func WriteIfStatmentToBuf(_buf *bytes.Buffer, ifExpr string, bodyFn func()) {
	if ifExpr == "true" {
		bodyFn()
		return
	} else if ifExpr == "false" {
		return
	}
	_buf.WriteString(`	if ` + ifExpr + `{
`)
	bodyFn()
	_buf.WriteString(`}
`)
}

package udwGoTypeMarshal

import "strings"

func StringToGoVarName(in string) string {
	in = strings.Replace(in, " ", "", -1)
	in = strings.Replace(in, ",", "", -1)
	in = strings.Replace(in, "/", "", -1)
	in = strings.Replace(in, "-", "", -1)
	in = strings.Replace(in, "'", "", -1)
	in = strings.Replace(in, ".", "", -1)
	return in
}

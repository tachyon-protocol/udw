package udwStrings

import (
	"strings"
)

func SplitLineTrimSpace(s string) []string {
	lines := strings.Split(s, "\n")
	output := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		output = append(output, line)
	}
	return output
}

func SplitLineNoEmptyLine(s string) []string {
	lines := strings.Split(s, "\n")
	output := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasSuffix(line, "\r") {
			line = line[:len(line)-1]
		}
		if line == "" {
			continue
		}
		output = append(output, line)
	}
	return output
}

func SplitLine(s string) []string {
	return strings.Split(s, "\n")
}
func TrimAllLineSpaceInPrefix(s string) string {
	tmp := strings.Split(strings.Replace(s, "\n", "", -1), " ")
	if len(tmp) == 0 {
		return s
	}
	var i int
	for {
		if i > len(tmp)-1 {
			break
		}
		if tmp[i] != "" {
			break
		}
		i++
	}
	return strings.Join(tmp[i:], " ")
}

func SplitTrimSpace(s string, sep string) []string {
	partList := strings.Split(s, sep)
	output := make([]string, 0, len(partList))
	for _, part := range partList {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		output = append(output, part)
	}
	return output
}

func GetSpaceArrayByIndex(s string, index int) string {
	parts := strings.Fields(s)
	if index < 0 || len(parts) >= index {
		return ""
	}
	return parts[index]
}

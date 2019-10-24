package udwSys

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"strings"
)

type IptableRule struct {
	Table string
	Rule  string
}

func MustSetIptableRule(rule IptableRule) {
	for _, thisRule := range MustGetIptableRuleList() {
		if thisRule.Table == rule.Table && thisRule.Rule == rule.Rule {
			return
		}
	}

	udwCmd.MustRun("iptables -w -t " + rule.Table + " " + rule.Rule)
}

func MustGetIptableRuleList() []IptableRule {
	content := udwCmd.MustCombinedOutput("iptables-save")
	return parseIptableSave(string(content))
}

func parseIptableSave(content string) []IptableRule {
	thisTable := ""
	output := []IptableRule{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == '#' || line[0] == ':' {
			continue
		}
		if line[0] == '*' {
			thisTable = line[1:]
			continue
		}
		if line == "COMMIT" {
			continue
		}
		output = append(output, IptableRule{
			Table: thisTable,
			Rule:  line,
		})
	}
	return output
}

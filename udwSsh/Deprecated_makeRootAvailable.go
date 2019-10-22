package udwSsh

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwStrings"
	"strings"
)

func MustMakeRootAvailableM1(remote RemoteServer) {
	SshCopyId(remote)

	runWithPassword := func(cmd string) []byte {
		if remote.Password != "" {
			return MustRpcSsh(remote, `echo `+remote.Password+` | sudo -S -- sh -c `+udwCmd.BashEscape(cmd))
		} else {
			return MustRpcSsh(remote, `sudo -S -- sh -c `+udwCmd.BashEscape(cmd))
		}
	}
	if remote.CertContent == "" {

		runWithPassword(`cp -r /home/` + remote.UserName + `/.ssh /root/`)
	} else {

		runWithPassword("mkdir -p /root/.ssh")
		runWithPassword("touch /root/.ssh/authorized_keys")
		content := runWithPassword("cat /root/.ssh/authorized_keys")
		newCertContent := addCertToauthorized_keysFile(string(content), string(remote.CertContent))
		runWithPassword("echo " + udwCmd.BashEscape(newCertContent) + " > /root/.ssh/authorized_keys")

	}
}

func addCertToauthorized_keysFile(oldContent string, toAddCert string) string {
	toAddCert = strings.TrimSpace(toAddCert)
	_buf := bytes.Buffer{}
	hasFound := false
	for _, line := range udwStrings.SplitLineTrimSpace(oldContent) {
		if strings.Contains(line, `command="`) {

			continue
		}
		_buf.WriteString(line)
		_buf.WriteByte('\n')
		if line == toAddCert {
			hasFound = true
		}
	}
	if !hasFound {
		_buf.WriteString(toAddCert)
		_buf.WriteByte('\n')
	}
	return _buf.String()
}

package udwSsh

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwFile"
	"path/filepath"
)

type MustSshInstallRootCertificationRequest struct {
	Ip                     string
	InPassword             string
	RootCertificateContent string
}

func MustSshInstallRootCertification(req MustSshInstallRootCertificationRequest) {
	udwCmd.MustRun("udw install sshpass")
	udwCmd.MustRunInBash(`ssh-keygen -f "` + filepath.Join(udwFile.MustGetHomeDirPath(), ".ssh/known_hosts") + `" -R ` + req.Ip)
	installCertS := `umask 077 && mkdir -p .ssh && echo ` + udwCmd.BashEscape(req.RootCertificateContent) + ` >> .ssh/authorized_keys`
	udwCmd.CmdBash(`sshpass -p ` + udwCmd.BashEscape(req.InPassword) + ` ssh -C -o StrictHostKeyChecking=no -o ConnectTimeout=10 root@` + req.Ip + ` ` +
		udwCmd.BashEscape(installCertS)).MustRun()
}

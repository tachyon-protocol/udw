package udwSsh

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwPlatform"
	"strconv"
)

type ScpReq struct {
	Ip                     string
	User                   string
	Port                   int
	IsUpload               bool
	LocalPath              string
	RemotePath             string
	OptionCompressionLevel int
}

func Scp(req ScpReq) error {
	if req.Port <= 0 {
		req.Port = 22
	}
	if req.User == "" {
		req.User = "root"
	}
	cmd := "scp -o StrictHostKeyChecking=no -o PasswordAuthentication=no -o ConnectTimeout=5 -o ServerAliveInterval=15 "
	if req.OptionCompressionLevel > 0 && udwPlatform.IsLinux() {
		if req.OptionCompressionLevel > 9 {
			req.OptionCompressionLevel = 9
		}
		cmd += "-o Compression=yes -o CompressionLevel=" + strconv.Itoa(req.OptionCompressionLevel) + " "
	}
	cmd += "-P " + strconv.Itoa(req.Port) + " "
	remotePath := req.User + "@" + req.Ip + ":" + req.RemotePath
	if req.IsUpload {
		cmd += req.LocalPath + " " + remotePath
	} else {
		cmd += remotePath + " " + req.LocalPath
	}
	return udwCmd.CmdBash(cmd).Run()
}

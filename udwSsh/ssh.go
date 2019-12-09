package udwSsh

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwRetry"
	"github.com/tachyon-protocol/udw/udwTime"
	"strconv"
	"strings"
	"time"
)

const SshOption = " -C -o StrictHostKeyChecking=no -o PasswordAuthentication=no "
const sshCmdPrefix = "ssh" + SshOption
const ScpCmdPrefix = "scp" + SshOption

type RemoteServer struct {
	Ip                string
	Port              int
	UserName          string
	ClientKeyFilePath string

	Password    string
	CertContent string

	IsScp         bool
	LocalPath     string
	RemotePath    string
	TimeoutSecond int
}

func (r RemoteServer) String() string {
	if r.Port == 0 {
		r.Port = 22
	}
	if r.UserName == "" {
		r.UserName = "root"
	}
	seg := []string{}
	portFlag := "-p"
	if r.ClientKeyFilePath != "" {
		seg = append(seg, "-i "+r.ClientKeyFilePath)
	}
	if r.IsScp {
		portFlag = "-P"
	}
	seg = append(seg, portFlag, strconv.Itoa(r.Port))
	if r.LocalPath != "" {
		seg = append(seg, r.LocalPath)
	}
	seg = append(seg, r.UserName+"@"+r.Ip)
	cmd := strings.Join(seg, " ")
	if r.RemotePath != "" {
		cmd += ":" + r.RemotePath
	}
	return cmd
}

func installSshpass() {
	udwCmd.MustRun("udw install sshpass")
}

func sshCopyIdWithSshpass(remote RemoteServer) {
	installSshpass()
	udwCmd.MustRunInBash(`sshpass -p "` + remote.Password + `" ssh-copy-id -o StrictHostKeyChecking=no ` + remote.String())
}

func SshCopyId(remote RemoteServer) {
	if remote.Ip == "" {
		return
	}
	isReachable, havePermission := availableCheckRemote(remote)
	if !isReachable {
		panic("[udwSsh SshCertCopyLocalToRemote]" + remote.String() + " unreachable!")
	}
	if havePermission {
		return
	}
	if remote.Password == "" {
		panic("password was not provided")
	}
	sshCopyIdWithSshpass(remote)
}

func MustSshCopyidWithIpAndCertContent(ip string, certContent string) {
	ret := MustRpcSshDefault(ip, `cat .ssh/authorized_keys`)
	if strings.Contains(string(ret), certContent) {

		return
	}
	ret = MustRpcSshDefault(ip, `umask 077 && mkdir -p .ssh && echo `+udwCmd.BashEscape(certContent)+` >> .ssh/authorized_keys && echo 'key add success'`)
	if string(ret) != "key add success\n" {
		panic("[MustSshCopyidWithIpAndCertContent] fail")
	}

}

func RpcSsh(remote RemoteServer, cmd ...string) (stdout []byte, err error) {
	if len(cmd) == 0 {
		return nil, err
	}
	if remote.Ip == "" {
		return nil, err
	}
	timeoutOption := ""
	if remote.TimeoutSecond > 0 {
		timeoutOption = "-o ConnectTimeout=" + strconv.Itoa(remote.TimeoutSecond) + " "
	}
	sshCmd := sshCmdPrefix + timeoutOption + remote.String() + " << EOF " + strings.Join(cmd, "&&") + "\nEOF"
	stdout, err = udwCmd.CmdString(sshCmd).RunAndReturnOutput()
	logPath := "/tmp/rpcSshCmd-" + remote.Ip
	udwFile.MustAppendFile(logPath, []byte(strings.Join([]string{sshCmd, udwTime.DefaultFormat(time.Now())}, "\n")))
	udwFile.MustAppendFile(logPath, stdout)
	if err != nil {
		errWithStdOut := fmt.Sprint(sshCmd, "\n", string(stdout), "\n", err)
		udwFile.MustAppendFile(logPath, []byte(errWithStdOut))
		return nil, err
	}
	return stdout, nil
}

func MustRpcSsh(remote RemoteServer, cmd ...string) []byte {
	out, err := RpcSsh(remote, cmd...)
	udwErr.PanicIfError(err)
	return out
}

func RpcSshDefault(ip string, cmd ...string) (b []byte, err error) {
	return RpcSsh(RemoteServer{Ip: ip}, cmd...)
}

func MustRpcSshDefault(ip string, cmd ...string) []byte {
	return MustRpcSsh(RemoteServer{Ip: ip}, cmd...)
}

func MustRpcSshDefaultWithBashContent(ip string, bash string) []byte {
	b, err := udwCmd.CmdBash(sshCmdPrefix + (&RemoteServer{Ip: ip}).String() + " bash -s").RunAndReturnOutputWithStdin([]byte(bash))
	udwErr.PanicIfError(err)
	return b
}

func AvailableCheck(ip string) (isReachable, havePermission bool) {
	return availableCheckRemote(RemoteServer{Ip: ip})
}

func availableCheckRemote(remote RemoteServer) (isReachable, havePermission bool) {
	err := udwRetry.Run(180, func() {
		logPrefix := "[udwSsh AvailableCheckRemote]"
		sign := "AvailableCheck__OK"
		cmdString := sshCmdPrefix + "-o ConnectTimeout=5 " + remote.String() + " echo " + sign
		cmd := udwCmd.CmdString(cmdString)
		b, e := cmd.CombinedOutput()
		if e == nil && strings.Contains(string(b), sign) {
			isReachable = true
			havePermission = true
			return
		}
		if e != nil && strings.Contains(string(b), "Permission denied") {
			isReachable = true
			havePermission = false
			return
		}
		panic(fmt.Sprintln(logPrefix, remote.String(), "UnReachable"))
	})
	if err == nil {
		return
	}
	fmt.Println(err)
	return false, false
}

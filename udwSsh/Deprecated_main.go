package udwSsh

import "github.com/tachyon-protocol/udw/udwCmd"

func MustScpToRemoteDefault(ip, localFileRealPath, remoteFileRealPath string) {
	MustScpToRemote(&RemoteServer{
		Ip:         ip,
		LocalPath:  localFileRealPath,
		RemotePath: remoteFileRealPath,
		IsScp:      true,
	})
}

func MustScpToRemote(remote *RemoteServer) {
	remote.IsScp = true
	udwCmd.CmdString(ScpCmdPrefix + remote.String()).SetDir("/").MustRunAndReturnOutput()
}

func ScpToRemoteDefault(ip, localFileRealPath, remoteFileRealPath string) error {
	return ScpToRemote(&RemoteServer{
		Ip:         ip,
		LocalPath:  localFileRealPath,
		RemotePath: remoteFileRealPath,
		IsScp:      true,
	})
}

func ScpToRemote(remote *RemoteServer) error {
	remote.IsScp = true
	return udwCmd.CmdString(ScpCmdPrefix + remote.String()).SetDir("/").Run()
}

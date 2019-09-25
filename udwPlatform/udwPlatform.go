package udwPlatform

import "runtime"

type Platform struct {
	Os   string
	Arch string
}

func (p Platform) Compatible(other Platform) bool {
	return p == other
}

func (p Platform) String() string {
	return p.Os + "_" + p.Arch
}

func (p Platform) GetExeSuffix() string {
	if p.Os == "windows" {
		return p.Os + "_" + p.Arch + ".exe"
	}
	return p.Os + "_" + p.Arch
}

var LinuxAmd64 = Platform{Os: "linux", Arch: "amd64"}
var DarwinAmd64 = Platform{Os: "darwin", Arch: "amd64"}
var WindowsAmd64 = Platform{Os: "windows", Arch: "amd64"}

func GetCompiledPlatform() Platform {
	return Platform{
		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsLinuxAmd64() bool {
	return runtime.GOOS == "linux" && runtime.GOARCH == "amd64"
}

func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

func IsDarwinAmd64() bool {
	return runtime.GOOS == "darwin" && runtime.GOARCH == "amd64"
}

func IsMac() bool {
	return GetCurrentUdwOs() == UdwOsMac
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsWindowsAmd64() bool {
	return runtime.GOOS == "windows" && runtime.GOARCH == "amd64"
}

func IsWindows386() bool {
	return runtime.GOOS == "windows" && runtime.GOARCH == "386"
}

func GoosStringIsWindows(goos string) bool {
	return goos == "windows"
}

func IsAndroid() bool {
	return runtime.GOOS == "android"
}

func IsIos() bool {
	return GetCurrentUdwOs() == string(UdwOsIos)
}

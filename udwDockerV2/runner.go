package udwDockerV2

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuild"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuildCtx"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"path/filepath"
	"strings"
)

func BuildRunnerToDownload(pkgPath string, os string) {
	resp := BuildRunner(pkgPath, os)
	copyFile := filepath.Join(udwFile.MustGetHomeDirPath(), "Downloads", getImageNameFromPkgPath(pkgPath)+resp.GetOutputExeFileExt())
	udwFile.MustCopy(resp.GetOutputExeFilePath(), copyFile)
}

func BuildRunner(pkgPath string, os string) (resp *udwGoBuildCtx.Ctx) {
	buildImage(pkgPath)
	imageFile := filepath.Join(getBuildPath(pkgPath), "image")
	exportImageToFile(pkgPath, imageFile)
	content := udwFile.MustReadFile(imageFile)
	buildPath := getBuildPath(pkgPath)
	runnerPath := filepath.Join(buildPath, "main.go")
	udwFile.MustMkdirForFile(runnerPath)
	udwFile.MustWriteFile(runnerPath, []byte(`package main
	
	import (
		"bytes"
		"github.com/tachyon-protocol/udw/udwCmd"
		"github.com/tachyon-protocol/udw/udwConsole"
		"github.com/tachyon-protocol/udw/udwErr"
		"os"
		"os/exec"
		"strconv"
		"fmt"
	)
	
	func main() {
		imageName := "`+getImageNameFromPkgPath(pkgPath)+`"
		udwConsole.AddCommandWithName("Run", func(req struct {
			Port    uint16
			Command string
		}) {
			cmd := exec.Command("docker", "image", "load")
			cmd.Stdin = bytes.NewReader([]byte(`+udwGoTypeMarshal.WriteStringToGolang(string(content))+`))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			udwErr.PanicIfError(err)
			err = udwCmd.Run("docker container rm --force " + imageName)
			if err != nil {
				fmt.Println(err)
			}
			command := "docker container run"
			if req.Port != 0 {
				command += " --publish " + strconv.Itoa(int(req.Port)) + ":" + strconv.Itoa(int(req.Port))
			}
			command += " --privileged --cap-add=NET_ADMIN --device=/dev/net/tun --name " + imageName + " " + imageName + " " + req.Command
			udwCmd.MustRun(command)
		})
		udwConsole.AddCommandWithName("Stop", func() {
			udwCmd.MustRun("docker container rm --force " + imageName)
		})
		udwConsole.Main()
	}`))
	runnerPkgPath := strings.TrimPrefix(
		strings.TrimPrefix(buildPath, filepath.Join(udwProjectPath.MustGetProjectPath(), "src")),
		"/",
	)
	resp = udwGoBuild.MustBuild(udwGoBuild.BuildRequest{
		PkgPath:       runnerPkgPath,
		TargetOs:      os,
		TargetCpuArch: `amd64`,
	})
	fmt.Println("- - -\nRunner:", resp.GetOutputExeFilePath())
	return resp
}

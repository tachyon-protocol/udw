package udwDockerV2

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuild"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"path/filepath"
	"strings"
)

func BuildImageToDownload(pkgPath string) {
	buildImage(pkgPath)
	exportImageToFile(pkgPath, filepath.Join(udwFile.MustGetHomeDirPath(), "Downloads", getImageNameFromPkgPath(pkgPath)+".image"))
}

func buildImage(pkgPath string) {
	resp := udwGoBuild.MustBuild(udwGoBuild.BuildRequest{
		PkgPath:       pkgPath,
		TargetOs:      `linux`,
		TargetCpuArch: `amd64`,
		EnableRace:    false,
	})
	packageName := filepath.Base(pkgPath)
	buildPath := getBuildPath(pkgPath)
	udwFile.MustMkdir(buildPath)
	udwFile.MustSetWd(buildPath)
	udwFile.MustCopy(resp.GetOutputExeFilePath(), filepath.Join(buildPath, packageName))

	dockerFile := `FROM ubuntu:18.04
WORKDIR /usr/local/bin
RUN apt-get update
RUN apt-get install -y net-tools
RUN apt-get install -y iptables
RUN apt-get install -y iproute2
COPY ` + packageName + ` ./
CMD ["` + packageName + `"]`
	udwFile.MustWriteFile(filepath.Join(buildPath, "Dockerfile"), []byte(dockerFile))
	udwCmd.MustRun(`docker image build -t ` + getImageNameFromPkgPath(pkgPath) + ` .`)
}

func exportImageToFile(pkgPath string, savePath string) {
	udwCmd.MustRun(`docker image save -o ` + savePath + ` ` + getImageNameFromPkgPath(pkgPath))
}

func getImageNameFromPkgPath(pkgPath string) (imageName string) {
	projectName := filepath.Base(udwProjectPath.MustGetProjectPath())
	imageName = strings.ToLower(strings.Join(append([]string{projectName}, strings.Split(pkgPath, "/")...), "-"))
	return imageName
}

func getBuildPath(pkgPath string) (buildPath string) {
	return udwProjectPath.MustPathInProject("src/tmp/dockerBuild/" + pkgPath)
}

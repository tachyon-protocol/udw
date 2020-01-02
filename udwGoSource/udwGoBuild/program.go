package udwGoBuild

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoBuildCtx"
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"github.com/tachyon-protocol/udw/udwStrings"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type programV2 struct {
	goroot                      string
	gopathList                  []string
	targetPackagePathOrFilePath string
	outputExeFilePath           string
	env                         map[string]string
	buildTagList                []string

	initLdFlags string
	dir         string
	ctx         *udwGoBuildCtx.Ctx
}

func mustNewProgramV2(env map[string]string, ctx *udwGoBuildCtx.Ctx) *programV2 {
	if env == nil {
		env = map[string]string{}
	}
	p := &programV2{
		env:                         env,
		targetPackagePathOrFilePath: ctx.BuildTargetPkgPath,
	}
	p.ctx = ctx

	p.goroot = p.ctx.GetGoRoot()
	p.gopathList = p.resolveGoPathList()
	p.gopathList = append(p.gopathList, p.goroot)
	p.initLdFlags = "-s -w"

	cgoCPPFLAGS := p.ctx.BuildCgoCppFlags
	cgoLDFLAGS := p.ctx.BuildCgoLdFlags
	for _, gopath := range p.gopathList {
		toAdd := "-I " + gopath
		if strings.Contains(cgoCPPFLAGS, toAdd) == false {
			cgoCPPFLAGS += " " + toAdd
		}
		if strings.Contains(cgoLDFLAGS, toAdd) == false {
			cgoLDFLAGS += " " + toAdd
		}

	}

	if p.ctx.IsGoOsDarwin() {

		cgoLDFLAGS += " -flat_namespace -Wl,-undefined,warning"
	} else {

		cgoLDFLAGS += " -Wl,--unresolved-symbols=ignore-all"
	}
	p.env["CGO_CPPFLAGS"] = cgoCPPFLAGS
	p.env["CGO_LDFLAGS"] = cgoLDFLAGS
	if p.ctx.BuildCgoCflags != "" {
		p.env["CGO_CFLAGS"] = p.ctx.BuildCgoCflags
	}
	return p
}

func mustNewProgramV2MergeDefaultEnv(env map[string]string, ctx *udwGoBuildCtx.Ctx) *programV2 {
	m := udwCmd.MustGetEnvMapFromSystem()
	if udwProjectPath.HasProjectPath() {
		m["GOPATH"] = udwProjectPath.MustGetProjectPath()
	}
	if env != nil {
		for k, v := range env {
			m[k] = v
		}
	}

	return mustNewProgramV2(m, ctx)
}

func (p *programV2) resolveGoRoot() string {
	return p.ctx.GetGoRoot()
}

func (p *programV2) resolveGoPathList() []string {
	return p.ctx.GetGoPathList()
}

func (p *programV2) GetGoos() string {
	return p.ctx.GetGoOs()
}

func (p *programV2) SetLdflags(s string) {
	p.initLdFlags = s
}

func (p *programV2) SetBuildTagList(BuildTagList []string) {
	p.buildTagList = BuildTagList
}

func (p *programV2) MustIsTargetExist() bool {
	if udwFile.MustOnlyFileExist(p.targetPackagePathOrFilePath) {
		return true
	}
	pkgDirPath := p.GetPackageFilePathByPackagePath(p.targetPackagePathOrFilePath)
	return pkgDirPath != ""
}

func (p *programV2) MustGoInstall() {
	p.mustUdwGoInstall(udwStrings.StringSliceMerge("go", "install", p.getBuildFlagCmdSlice(), p.targetPackagePathOrFilePath))
}

func (p *programV2) getBuildFlagCmdSlice() (output []string) {
	if p.ctx.EnableRace {
		output = udwStrings.StringSliceMerge(output, "-race")
	}
	if len(p.buildTagList) > 0 {
		output = udwStrings.StringSliceMerge(output, "-tags", strings.Join(p.buildTagList, " "))
	}
	if p.ctx.BuildCgoDebug {
		output = udwStrings.StringSliceMerge(output, "-n")
	}
	ldflagsBuf := &bytes.Buffer{}
	ldflagsBuf.WriteString(p.initLdFlags)
	ldflagsBuf.WriteString(" ")
	variableMap := p.ctx.BuildVariableMap
	if variableMap != nil && len(variableMap) > 0 {
		pairList := udwMap.MapStringStringToKeyValuePairListAes(variableMap)
		for _, pair := range pairList {
			ldflagsBuf.WriteString("-X ")
			ldflagsBuf.WriteString(pair.Key)
			ldflagsBuf.WriteString("=")
			ldflagsBuf.WriteString(pair.Value)
			ldflagsBuf.WriteString(" ")
		}
	}
	s := strings.TrimSpace(ldflagsBuf.String())
	if s != "" {
		output = udwStrings.StringSliceMerge(output, "-ldflags", s)
	}
	if len(p.gopathList) > 0 {
		output = udwStrings.StringSliceMerge(output, "-gcflags=-trimpath="+p.gopathList[0])
	}
	return output
}

func (p *programV2) GetPackageFilePathByPackagePath(PackagePath string) string {
	for _, gopath := range p.gopathList {
		thisPath := filepath.Join(gopath, "src", PackagePath)
		fi, err := os.Stat(thisPath)
		if err != nil {
			if udwFile.ErrorIsFileNotFound(err) {
				continue
			}
			panic(err)
		}
		if !fi.IsDir() {
			continue
		}
		return thisPath
	}
	return ""
}

func (p *programV2) GetGoArch() string {
	goarch := p.env["GOARCH"]
	if goarch == "" {
		return runtime.GOARCH
	}
	return goarch
}

func (p *programV2) mustUdwGoInstall(cmdSlice []string) {
	for i := 0; i < 10; i++ {

		output, err := udwCmd.CmdSlice(cmdSlice).
			MustSetEnvMap(p.env).
			CombinedOutput()
		if err == nil {
			fmt.Println("GOROOT="+p.goroot,
				"GOPATH="+strings.Join(p.gopathList, ":"),
				strings.Join(cmdSlice, " "))
			fmt.Print(string(output))
			return
		}

		outputS := string(output)
		if strings.Contains(outputS, "permission denied") {
			udwCmd.MustRun("sudo chmod -R 777 " + filepath.Join(p.goroot, "pkg"))
			continue
		}
		if runtime.GOOS == `darwin` && strings.Contains(outputS, "xcrun: error: invalid active developer path") {

			udwCmd.CmdSlice([]string{"xcode-select", "--install"}).MustRun()
			continue
		}
		hasModify := fixImportAndNotUsedError(outputS)
		if hasModify {
			continue
		}
		errMsg := "mustUdwGoInstall 2ern9u7k6h GoOs:" + p.GetGoos() +
			" GoArch:" + p.GetGoArch() +
			" Cmd:" + strings.Join(cmdSlice, " ") + "\n" +
			"GOROOT=" + p.goroot +
			" GOPATH=" + strings.Join(p.gopathList, ":") + "\n" +
			"CmdIutput:" + string(output) + "\nErrMsg:" + udwErr.ErrorToMsg(err)
		fmt.Println(errMsg)
		panic(errMsg)
	}
}

func fixImportAndNotUsedError(outputS string) bool {
	importedAndNotUsedS := "imported and not used: "
	redeclaredS := "redeclared as imported package name"
	hasModify := false
	for _, line := range udwStrings.SplitLineTrimSpace(outputS) {
		if strings.Contains(line, importedAndNotUsedS) {

			beforePart := line[:strings.Index(line, importedAndNotUsedS)]
			beforePartList := strings.Split(beforePart, ":")
			if len(beforePartList) != 3 && len(beforePartList) != 4 {
				fmt.Println("error", "[fixImportAndNotUsedError] importedAndNotUsedS len(beforePartList)!=3", len(beforePartList), line)
				continue
			}
			path := strings.TrimSpace(beforePartList[0])
			lineNum, err := strconv.Atoi(beforePartList[1])
			if err != nil {
				fmt.Println("error", "[fixImportAndNotUsedError] importedAndNotUsedS strconv.Atoi(beforePartList[1])", beforePartList[1], line)
				continue
			}
			packageName := strings.TrimSpace(line[strings.Index(line, importedAndNotUsedS)+len(importedAndNotUsedS):])
			if !udwFile.MustFileExist(path) {
				fmt.Println("error", "[fixImportAndNotUsedError] !udwFile.MustFileExist(path)", path, line)
				continue
			}
			err = commentOutCodeByPathAndLineNum(path, lineNum, packageName)
			if err != nil {
				fmt.Println("error", "[fixImportAndNotUsedError]", err.Error(), line)
				continue
			}
			fmt.Println("fix import and not use error 1", path, lineNum, line)
			hasModify = true
		} else if strings.Contains(line, redeclaredS) {

			beforePart := line[:strings.Index(line, redeclaredS)]
			beforePartList := strings.Split(beforePart, ":")
			if len(beforePartList) != 3 && len(beforePartList) != 4 {
				fmt.Println("error", "[fixImportAndNotUsedError] redeclaredS len(beforePartList)!=3", len(beforePartList), line)
				continue
			}
			path := strings.TrimSpace(beforePartList[0])
			lineNum, err := strconv.Atoi(beforePartList[1])
			if err != nil {
				fmt.Println("error", "[fixImportAndNotUsedError] redeclaredS strconv.Atoi(beforePartList[1])", beforePartList[1], line)
				continue
			}
			packageName := strings.TrimSpace(beforePartList[2])
			err = commentOutCodeByPathAndLineNum(path, lineNum, packageName)
			if err != nil {
				fmt.Println("error", "[fixImportAndNotUsedError]", err.Error(), line)
				continue
			}
			fmt.Println("fix import and not use error 2", path, lineNum, line)
			hasModify = true
		}
	}
	return hasModify
}

func commentOutCodeByPathAndLineNum(path string, lineNum int, shouldContains string) (err error) {
	if !udwFile.MustFileExist(path) {
		return udwErr.ErrorSprint("[commentOutCodeByPathAndLineNum] !udwFile.MustFileExist(path)", path, lineNum)
	}
	writeBackBuf := bytes.Buffer{}
	fileContent := udwFile.MustReadFile(path)
	linePart := udwStrings.SplitLine(string(fileContent))
	for i, line := range linePart {
		if i == lineNum-1 {
			if !strings.Contains(line, shouldContains) {
				return udwErr.ErrorSprint("[commentOutCodeByPathAndLineNum] !strings.Contains(line,shouldContains)", line, shouldContains)
			}
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				return nil
			}
			writeBackBuf.WriteString("//")
		}
		writeBackBuf.WriteString(line)
		if i != len(linePart)-1 {
			writeBackBuf.WriteByte('\n')
		}
	}
	udwFile.MustWriteFile(path, writeBackBuf.Bytes())
	return nil
}

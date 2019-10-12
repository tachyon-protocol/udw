package udwGoWriter

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/kgwtPkg"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"github.com/tachyon-protocol/udw/udwTest"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteFileGetter(ot *testing.T) {
	udwFile.MustDelete("zzzig_test")
	defer udwFile.MustDelete("zzzig_test")
	testWriteFileGetterL1()
	testWriteFileGetterL2()
	testWriteFileGetterL3()
}

func testWriteFileGetterL1() {
	toWriteData := make([]byte, 65536)
	for i := 0; i < 65536; i += 2 {
		toWriteData[i] = byte(i % 256)
		toWriteData[i+1] = byte(i / 255)
	}
	rootPath := udwProjectPath.MustGetProjectPath()
	WriteGoFileGetterWithGlobalVariable(WriteGoFileGetterRequest{
		PackageName:  "main",
		FunctionName: "getData",
		Obj:          toWriteData,
		GoFilePath:   filepath.Join(rootPath, "src/github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/zzzig_test/kgwt01/getData.go"),
	})
	udwFile.MustWriteFileWithMkdir(filepath.Join(rootPath, "src/github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/zzzig_test/kgwt01/main.go"), []byte(`package main

import "fmt"

func main(){
	fmt.Println(getData())
}`))
	output := udwCmd.CmdString("go run github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/zzzig_test/kgwt01").MustSetEnv("GOPATH", rootPath).MustCombinedOutput()
	udwTest.Equal(output, []byte(fmt.Sprintln(toWriteData)))
}

func testWriteFileGetterL2() {
	v1 := &kgwtPkg.T1{
		M: map[string]string{
			"a": "a",
		},
		Now: time.Now(),
	}
	rootPath := udwProjectPath.MustGetProjectPath()
	WriteGoFileGetterWithGlobalVariable(WriteGoFileGetterRequest{
		PackageName:  "main",
		FunctionName: "getData",
		Obj:          v1,
		GoFilePath:   filepath.Join(rootPath, "src/github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/zzzig_test/kgwt02/getData.go"),
	})
	udwFile.MustWriteFileWithMkdir(filepath.Join(rootPath, "src/github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/zzzig_test/kgwt02/main.go"), []byte(`package main

import "os"

func main(){
	os.Stdout.WriteString(getData().String());
}`))
	output, _ := udwCmd.CmdString("go run github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/zzzig_test/kgwt02").MustSetEnv("GOPATH", rootPath).CombinedOutput()
	udwTest.Equal(output, []byte(v1.String()))
}

func testWriteFileGetterL3() {
	v1 := time.Now().UTC()
	rootPath := udwProjectPath.MustGetProjectPath()
	WriteGoFileGetterWithGlobalVariable(WriteGoFileGetterRequest{
		PackageName:  "main",
		FunctionName: "getData",
		Obj:          v1,
		GoFilePath:   filepath.Join(rootPath, "src/github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/zzzig_test/kgwt03/getData.go"),
	})
	udwFile.MustWriteFileWithMkdir(filepath.Join(rootPath, "src/github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/zzzig_test/kgwt03/main.go"), []byte(`package main

import "os"

func main(){
	os.Stdout.WriteString(getData().String());
}`))
	output, _ := udwCmd.CmdString("go run github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/zzzig_test/kgwt03").MustSetEnv("GOPATH", rootPath).CombinedOutput()
	udwTest.Equal(output, []byte(v1.String()))
}

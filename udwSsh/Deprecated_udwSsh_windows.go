// +build windows

package udwSsh

import (
	"os"
	osUser "os/user"
	"path/filepath"
	"strings"
)

func certFilePathList() []string {
	u, err := osUser.Current()
	if err != nil {
		return nil
	}

	output := []string{}
	output = append(output, certFilePathListFromCygwin(getOsUserName(u))...)
	output = append(output, filepath.Join(u.HomeDir, ".ssh", "id_rsa"))
	return output
}
func certFilePathListFromCygwin(username string) []string {
	paths := getCygwinInstallPaths()
	output := make([]string, len(paths))
	for i, path := range paths {
		output[i] = filepath.Join(path, "home", username, ".ssh", "id_rsa")
	}
	return output
}

func getCygwinInstallPaths() []string {
	pathenv := os.Getenv("PATH")
	if pathenv == "" {
		return nil
	}
	output := []string{}
	for _, dir := range strings.Split(pathenv, ":") {
		if strings.Index(strings.ToLower(dir), "cygwin") >= 0 {
			output = append(output, filepath.Dir(dir))
		}
	}
	return output
}

func getOsUserName(u *osUser.User) string {
	pos := strings.LastIndex(u.Username, "\\")
	return u.Username[pos+1:]
}

package udwFile

import (
	"os/user"
)

func MustGetHomeDirPath() string {
	path, err := GetHomeDirPath()
	if err != nil {
		panic(err)
	}
	return path
}

func GetHomeDirPath() (string, error) {

	_user, err := user.Current()
	if err != nil {
		return "", err
	}
	return _user.HomeDir, nil
}

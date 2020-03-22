package pathmgr

import (
	"errors"
	"os"
	"runtime"
)

type EnvPath struct {
	homeDir *string
	cwd     *string
}

var envPath *EnvPath

func HomeDir() (string, error) {
	if envPath.homeDir == nil {
		var envName string
		switch runtime.GOOS {
		case "windows":
			envName = "USERPROFILE"
		default:
			envName = "HOME"
		}
		if hdTmp := os.Getenv(envName); len(hdTmp) == 0 {
			return "", errors.New("Home directory is not set in environment variables.")
		} else {
			envPath.homeDir = &hdTmp
		}
	}
	return *(envPath.homeDir), nil
}

func CWD() (string, error) {
	if envPath.cwd == nil {
		if cwdTmp, err := os.Getwd(); err == nil {
			envPath.cwd = &cwdTmp
		} else {
			return "", nil
		}
	}
	return *(envPath.cwd), nil
}

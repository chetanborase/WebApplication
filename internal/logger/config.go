package logger

import (
	"WebApplication/env"
	"path/filepath"
)

func infoFilePath() string {
	return filepath.Join(env.Get(env.LogInfoDirName, "logfiles"), "info.log")
}
func errFilePath() string {
	return filepath.Join(env.Get(env.LogInfoDirName, "logfiles"), "info.log")
}

//max age in days
func maxAgeOfLogFile() int {
	return 30
}

func maxBackup() int {
	return 0
}

//max size in MB
func maxSize() int {
	return 2
}

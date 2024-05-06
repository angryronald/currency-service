package file

import (
	"errors"
	"os"
)

var ErrProjectDirectoryNotFound = errors.New("project directory not found")

func firstIndexOfString(source string, value string) int {
	for i := 0; i+len(value) <= len(source); i++ {
		if source[i:i+len(value)] == value {
			return i
		}
	}
	return -1
}

func GetProjectDir(projectName string) (string, error) {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	index := firstIndexOfString(currentWorkingDir, projectName)
	if index == -1 {
		return "", ErrProjectDirectoryNotFound
	}

	return currentWorkingDir[:index+len(projectName)], nil
}

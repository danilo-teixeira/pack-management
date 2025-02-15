package helpers

import (
	"os"
	"strings"
)

func GetRootDirectory() (string, error) {
	dir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dirs := strings.Split(dir, "/")
	rootDir := ""

	for _, dir := range dirs {
		rootDir += dir + "/"

		info, _ := os.Stat(rootDir + "go.mod")

		if info != nil {
			if info.Name() == "go.mod" {
				break
			}
		}
	}

	return rootDir, nil
}

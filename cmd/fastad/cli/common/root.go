package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func GetFastADRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current working directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(cwd, ".fastad_root")); err == nil {
			return cwd, nil
		}

		if cwd == "/" {
			return "", errors.New("fastad root not found")
		}

		cwd = filepath.Dir(cwd)
	}
}

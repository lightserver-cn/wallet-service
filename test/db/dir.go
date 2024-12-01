package db

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetDirPath() (string, error) {
	if os.Getenv("MOCK_CALLER_FAILURE") == "true" {
		return "", fmt.Errorf("failed to get the current file path")
	}

	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		return "", fmt.Errorf("failed to get the current file path")
	}

	return filepath.Dir(filename), nil
}

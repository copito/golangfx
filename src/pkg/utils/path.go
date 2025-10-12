package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetModuleRoot() (string, error) {
	const moduleRootMarker = "go.mod"

	_, err := os.ReadFile(moduleRootMarker)
	if err != nil {
		// If go.mod is not in the current directory, bubble up until you find
		currentDir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}

		for {
			modFilePath := filepath.Join(currentDir, moduleRootMarker)
			if _, err := os.Stat(modFilePath); err == nil {
				return currentDir, nil
			}
			parentDir := filepath.Dir(currentDir)
			if parentDir == currentDir {
				return "", fmt.Errorf("%v not found in current or parent directories", moduleRootMarker)
			}
			currentDir = parentDir
		}
	}
	return filepath.Dir(filepath.Join(os.Getenv("PWD"), moduleRootMarker)), nil
}

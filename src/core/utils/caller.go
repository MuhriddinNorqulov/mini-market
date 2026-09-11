package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func CallerPath(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	absPath, err := filepath.Abs(file)
	if err != nil {
		absPath = file
	}
	link := fmt.Sprintf("file://%s:%d", absPath, line)
	return link
}

package utils

import (
	"os"
	"time"
	"image_processor-go/internal/core"
)

func GetPriorityByAge(fileModTime time.Time) core.Priority {
	age := time.Since(fileModTime)

	switch {
	case age > 365*24*time.Hour:
		return core.PriorityHigh
	case age > 30*24*time.Hour:
		return core.PriorityNormal
	default:
		return core.PriorityLow
	}
}

func GetFileInfo(path string) (modTime time.Time, size int64, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, 0, err
	}

	return info.ModTime(), info.Size(), nil
}

package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/disintegration/imaging"
	"image_processor-go/internal/core"
)

func ConvertImage(job *core.Job) (time.Duration, error) {
	start := time.Now()

	src, err := imaging.Open(job.InputPath)
	if err != nil {
		return 0, fmt.Errorf("open image: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(job.OutputPath), 0755); err != nil {
		return 0, fmt.Errorf("create output dir: %w", err)
	}

	switch job.TargetFormat {
	case core.FormatPNG:
		err = imaging.Save(src, job.OutputPath)
	case core.FormatWEBP:
		err = imaging.Save(src, job.OutputPath, imaging.JPEGQuality(job.Quality))
	case core.FormatJPEG:
		fallthrough
	default:
		err = imaging.Save(src, job.OutputPath, imaging.JPEGQuality(job.Quality))
	}

	if err != nil {
		return 0, fmt.Errorf("Save image: %w", err)
	}
	
	return time.Since(start), nil
}

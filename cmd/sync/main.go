package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"image_processor-go/internal/core"
	"image_processor-go/internal/generator"
	"image_processor-go/internal/processor"
)

func main() {
	start := time.Now()

	imageCount := 100
	inputDir := "../../test_input_sync"
	outputDir := "../../test_output_sync"
	targetFormat := core.FormatPNG

	fmt.Println("Generating test images...")
	if err := generator.GenerateTestImages(imageCount, inputDir); err != nil {
		fmt.Printf("Error generating images: %v\n", err)
		return
	}

	files, err := os.ReadDir(inputDir)
    if err != nil {
        fmt.Printf("Error reading dir: %v\n", err)
        return
    }

	fmt.Printf("Processing %d images synchronously...\n", len(files))

	processed := 0
	for i, file := range files {
		if file.IsDir() {
			continue
		}

		inputPath := filepath.Join(inputDir, file.Name())
		outputName := fmt.Sprintf("img_%04d.%s", i+1, targetFormat)
		outputPath := filepath.Join(outputDir, outputName)

		info, err := os.Stat(inputPath)
		if err != nil {
			fmt.Printf("Error stating %s: %v\n", inputPath, err)
			continue
		}

		job := &core.Job{
            ID:          i + 1,
            InputPath:   inputPath,
            OutputPath:  outputPath,
            FileModTime: info.ModTime(),
            FileSize:    info.Size(),
            TargetFormat: targetFormat,
            Quality:     85,
        }

		age := time.Since(info.ModTime())
		switch {
		case age > 375*24*time.Hour:
			job.Priority = core.PriorityHigh
		case age > 30*24*time.Hour:
			job.Priority = core.PriorityMid
		default:
			job.Priority = core.PriorityLow
		}

		duration, err := processor.ConvertImage(job)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", inputPath, err)
		} else {
			fmt.Printf(
				"[%d/%d] %s -> %s (priority: %d, time: %v)\n",
				i+1,
				len(files),
				file.Name(),
				outputName,
				job.Priority,
				duration,
			)
			processed++
		}

		if (i+1) % 10 == 0 {
			fmt.Printf("Progress: %d/%d\n", i+1, len(files))
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("\n=== RESULTS ===\n")
    fmt.Printf("Total time: %v\n", elapsed)
    fmt.Printf("Processed: %d/%d\n", processed, len(files))
    fmt.Printf("Speed: %.2f images/second\n", float64(processed)/elapsed.Seconds())
}

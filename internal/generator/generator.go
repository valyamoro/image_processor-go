package generator

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"path/filepath"
	"time"
	"github.com/disintegration/imaging"
)

func GenerateTestImages(count int, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	fmt.Printf("Generating %d test images in %s...\n", count, outputDir)

	sizes := []struct{ width, height int }{
		{800, 600}, {1024, 768}, {640, 480},
		{1280, 720}, {1920, 1080},
	}

	for i := 1; i <= count; i++ {
		size := sizes[rand.Intn(len(sizes))]

		img := imaging.New(size.width, size.height, color.White)

		for j := 0; j < 3; j++ {
			x := rand.Intn(size.width - 100)
			y := rand.Intn(size.height - 100)
			width := rand.Intn(100) + 20
			height := rand.Intn(100) + 20
			col := color.RGBA{
				R: uint8(rand.Intn(200)),
				G: uint8(rand.Intn(200)),
				B: uint8(rand.Intn(200)),
				A: 255,
			}
			
			rect := imaging.New(width, height, col)
			img = imaging.Overlay(img, rect, image.Pt(x, y), 1.0)
		}

		filename := filepath.Join(outputDir, fmt.Sprintf("img_%04d.jpg", i))
		quality := rand.Intn(20) + 80
		if err := imaging.Save(img, filename, imaging.JPEGQuality(quality)); err != nil {
			return fmt.Errorf("failed to save %s: %w", filename, err)
		}

		daysAgo := rand.Intn(730) + 1
		modTime := time.Now().Add(-time.Duration(daysAgo) * 24 * time.Hour)
		if err := os.Chtimes(filename, modTime, modTime); err != nil {
			return fmt.Errorf("failed to set time for %s: %w", filename, err)
		}

		if i%100 == 0 {
			fmt.Printf("Created %d/%d images\n", i, count)
		}
	}

	fmt.Println("Generation completed!")
	return nil
}

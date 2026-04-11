package converter

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a simple 100x100 test image in memory
func createDummyImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Fill the left half with white and the right half with black
	draw.Draw(img, image.Rect(0, 0, 50, 100), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(50, 0, 100, 100), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	return img
}

func TestProcessAndSave_Integration(t *testing.T) {
	// 1. Setup: Create a temporary directory that Go will automatically clean up after the test!
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test_output.png")

	// 2. Setup: Create our dummy image
	dummyImg := createDummyImage()

	// 3. Setup: Define a standard configuration
	cfg := Config{
		Width:            50,
		FontSize:         8.0,
		DPI:              72.0,      // Use a lower DPI so the test runs blazingly fast
		SkipHex:          "#000000", // Try skipping the black half of our dummy image!
		TolerancePercent: 2.0,
		OutputPath:       outputPath,
	}

	// 4. Execution: Run the actual conversion engine
	err := ProcessAndSave(dummyImg, cfg)
	// 5. Assertion: Did it return an error?
	if err != nil {
		t.Fatalf("ProcessAndSave failed unexpectedly: %v", err)
	}

	// 6. Assertion: Was the file actually created on the hard drive?
	fileInfo, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		t.Fatalf("Expected output file was not created at %s", outputPath)
	}

	// 7. Assertion: Does the file have actual content (not 0 bytes)?
	if fileInfo.Size() == 0 {
		t.Errorf("Output file was created but is completely empty (0 bytes)")
	}
}

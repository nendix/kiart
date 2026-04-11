package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to create a physical test image on the hard drive
func createPhysicalTestImage(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func TestCLI_EndToEnd(t *testing.T) {
	// 1. Setup a temporary directory for our test files
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.png")
	outputPath := filepath.Join(tempDir, "output.png")

	// 2. Create a dummy image to convert
	err := createPhysicalTestImage(inputPath)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectExitCode int
		expectOutput   string
	}{
		{
			name:           "Missing arguments (Fail)",
			args:           []string{"main.go"}, // No image provided
			expectExitCode: 1,
			expectOutput:   "Usage:", // Should print the help menu
		},
		{
			name:           "Invalid input file (Fail)",
			args:           []string{"main.go", "does_not_exist.jpg"},
			expectExitCode: 1,
			expectOutput:   "Error opening image",
		},
		{
			name: "Successful conversion (Pass)",
			// Note: We use short flags and small width/size to make the test fast
			args:           []string{"main.go", "-w", "50", "-s", "5", "-o", outputPath, inputPath},
			expectExitCode: 0,
			expectOutput:   "Successfully generated ASCII art",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create the command: `go run main.go [args...]`
			// Note: We prepend "run" to the args list
			cmdArgs := append([]string{"run"}, tc.args...)
			cmd := exec.Command("go", cmdArgs...)

			// Capture both stdout and stderr
			outputBytes, err := cmd.CombinedOutput()
			outputStr := string(outputBytes)

			// Check Exit Codes
			// If err is nil, the exit code is 0. If it's an ExitError, we can check the code.
			actualExitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					actualExitCode = exitErr.ExitCode()
				} else {
					t.Fatalf("Failed to run command: %v", err)
				}
			}

			if actualExitCode != tc.expectExitCode {
				t.Errorf("Expected exit code %d, got %d. Output: %s", tc.expectExitCode, actualExitCode, outputStr)
			}

			// Check if the expected text was printed to the terminal
			if !strings.Contains(outputStr, tc.expectOutput) {
				t.Errorf("Expected terminal output to contain '%s', but got: \n%s", tc.expectOutput, outputStr)
			}
		})
	}

	// Final verification for the successful run: Does the output file actually exist?
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("The success test passed, but the output file %s was never actually created on disk!", outputPath)
	}
}

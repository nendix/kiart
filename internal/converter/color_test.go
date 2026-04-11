package converter

import (
	"image/color"
	"math"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		name        string
		hexStr      string
		expected    color.RGBA
		expectError bool
	}{
		{"Valid White with hash", "#FFFFFF", color.RGBA{255, 255, 255, 255}, false},
		{"Valid Black without hash", "000000", color.RGBA{0, 0, 0, 255}, false},
		{"Valid Red", "#FF0000", color.RGBA{255, 0, 0, 255}, false},
		{"Invalid short hex", "#FFF", color.RGBA{}, true},
		{"Invalid characters", "#GGGGGG", color.RGBA{}, true},
		{"Empty string", "", color.RGBA{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseHexColor(tc.hexStr)

			if tc.expectError && err == nil {
				t.Errorf("Expected an error for input '%s', but got none", tc.hexStr)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Did not expect an error for '%s', but got: %v", tc.hexStr, err)
			}
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestColorDistance(t *testing.T) {
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	red := color.RGBA{255, 0, 0, 255}

	// The maximum distance is sqrt(255^2 + 255^2 + 255^2)
	maxDist := math.Sqrt(255*255 + 255*255 + 255*255)

	tests := []struct {
		name     string
		c1       color.RGBA
		c2       color.RGBA
		expected float64
	}{
		{"Exact same color (Black)", black, black, 0.0},
		{"Exact same color (Red)", red, red, 0.0},
		{"Absolute opposites (Black and White)", black, white, maxDist},
		{"Red to Black", red, black, 255.0}, // sqrt(255^2 + 0 + 0) = 255
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := colorDistance(tc.c1, tc.c2)
			// Using a small tolerance for floating point comparison
			if math.Abs(result-tc.expected) > 0.001 {
				t.Errorf("Expected distance %f, got %f", tc.expected, result)
			}
		})
	}
}

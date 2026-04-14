package color

import (
	"image/color"
	"math"
	"testing"
)

func TestParseHex(t *testing.T) {
	tests := []struct {
		name        string
		hexStr      string
		expected    color.RGBA
		expectError bool
	}{
		{"Valid White with hash", "#FFFFFF", color.RGBA{255, 255, 255, 255}, false},
		{"Valid Black without hash", "000000", color.RGBA{0, 0, 0, 255}, false},
		{"Valid Red", "#FF0000", color.RGBA{255, 0, 0, 255}, false},
		{"Transparent", "transparent", color.RGBA{0, 0, 0, 0}, false},
		{"Invalid short hex", "#FFF", color.RGBA{}, true},
		{"Invalid characters", "#GGGGGG", color.RGBA{}, true},
		{"Empty string", "", color.RGBA{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseHex(tc.hexStr)

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

func TestParseHexWithDefault(t *testing.T) {
	defaultColor := color.RGBA{128, 128, 128, 255}

	result, err := ParseHexWithDefault("", defaultColor)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != defaultColor {
		t.Errorf("Expected default %v, got %v", defaultColor, result)
	}

	result, err = ParseHexWithDefault("#FF0000", defaultColor)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != (color.RGBA{255, 0, 0, 255}) {
		t.Errorf("Expected red, got %v", result)
	}
}

func TestDistance(t *testing.T) {
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	red := color.RGBA{255, 0, 0, 255}

	maxDist := math.Sqrt(255*255 + 255*255 + 255*255)

	tests := []struct {
		name     string
		c1       color.RGBA
		c2       color.RGBA
		expected float64
	}{
		{"Same color (Black)", black, black, 0.0},
		{"Same color (Red)", red, red, 0.0},
		{"Opposites (Black and White)", black, white, maxDist},
		{"Red to Black", red, black, 255.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Distance(tc.c1, tc.c2)
			if math.Abs(result-tc.expected) > 0.001 {
				t.Errorf("Expected distance %f, got %f", tc.expected, result)
			}
		})
	}
}

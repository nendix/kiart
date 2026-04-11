package converter

import (
	"fmt"
	"image/color"
	"math"
	"strings"
)

func ParseHexColor(hex string) (color.RGBA, error) {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b uint8
	if len(hex) == 6 {
		// Capture the number of parsed items (n) and any error
		n, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)

		// If it failed to parse exactly 3 items, or threw an error, reject it
		if err != nil || n != 3 {
			return color.RGBA{}, fmt.Errorf("invalid hex characters: %s", hex)
		}

		return color.RGBA{r, g, b, 255}, nil
	}
	return color.RGBA{}, fmt.Errorf("invalid hex format, use #RRGGBB")
}

// colorDistance calculates the 3D distance between two colors
func colorDistance(c1, c2 color.RGBA) float64 {
	rDiff := float64(c1.R) - float64(c2.R)
	gDiff := float64(c1.G) - float64(c2.G)
	bDiff := float64(c1.B) - float64(c2.B)
	return math.Sqrt(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff)
}

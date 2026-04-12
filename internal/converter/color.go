package converter

import (
	"fmt"
	"image/color"
	"math"
	"strings"
)

func ParseHexColor(hex string) (color.RGBA, error) {
	if strings.ToLower(hex) == "transparent" {
		return color.RGBA{0, 0, 0, 0}, nil
	}

	hex = strings.TrimPrefix(hex, "#")
	var r, g, b, a uint8

	if len(hex) == 6 {
		n, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		if err != nil || n != 3 {
			return color.RGBA{}, fmt.Errorf("invalid hex characters: %s", hex)
		}
		return color.RGBA{r, g, b, 255}, nil
	}

	if len(hex) == 8 {
		n, err := fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
		if err != nil || n != 4 {
			return color.RGBA{}, fmt.Errorf("invalid hex characters: %s", hex)
		}

		// GO GOTCHA: color.RGBA requires pre-multiplied alpha!
		// We must mathematically reduce the RGB values based on the Alpha value.
		rr := uint8((uint16(r) * uint16(a)) / 255)
		gg := uint8((uint16(g) * uint16(a)) / 255)
		bb := uint8((uint16(b) * uint16(a)) / 255)

		return color.RGBA{rr, gg, bb, a}, nil
	}

	return color.RGBA{}, fmt.Errorf("invalid hex format, use #RRGGBB or #RRGGBBAA")
}

// colorDistance calculates the 3D distance between two colors
func colorDistance(c1, c2 color.RGBA) float64 {
	rDiff := float64(c1.R) - float64(c2.R)
	gDiff := float64(c1.G) - float64(c2.G)
	bDiff := float64(c1.B) - float64(c2.B)
	return math.Sqrt(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff)
}

package color

import (
	"fmt"
	"image/color"
	"math"
	"strings"
)

// ParseHex parses a hex color string (#RRGGBB or #RRGGBBAA) into a color.RGBA.
// Also accepts "transparent" as a special value.
func ParseHex(hex string) (color.RGBA, error) {
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
		// color.RGBA requires pre-multiplied alpha
		rr := uint8((uint16(r) * uint16(a)) / 255)
		gg := uint8((uint16(g) * uint16(a)) / 255)
		bb := uint8((uint16(b) * uint16(a)) / 255)
		return color.RGBA{rr, gg, bb, a}, nil
	}

	return color.RGBA{}, fmt.Errorf("invalid hex format, use #RRGGBB or #RRGGBBAA")
}

// ParseHexWithDefault parses a hex color string, returning defaultColor if hex is empty.
func ParseHexWithDefault(hex string, defaultColor color.RGBA) (color.RGBA, error) {
	if hex == "" {
		return defaultColor, nil
	}
	return ParseHex(hex)
}

// Distance calculates the Euclidean distance between two colors in RGB space.
func Distance(c1, c2 color.RGBA) float64 {
	rDiff := float64(c1.R) - float64(c2.R)
	gDiff := float64(c1.G) - float64(c2.G)
	bDiff := float64(c1.B) - float64(c2.B)
	return math.Sqrt(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff)
}

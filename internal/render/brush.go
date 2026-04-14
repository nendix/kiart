package render

import (
	"image"
	"image/color"
)

// buildShadedLUT pre-computes 256 shaded color variations for a base color.
func buildShadedLUT(baseColor color.RGBA) [256]*image.Uniform {
	var lut [256]*image.Uniform
	for i := range 256 {
		r := uint8((uint16(baseColor.R) * uint16(i)) / 255)
		g := uint8((uint16(baseColor.G) * uint16(i)) / 255)
		b := uint8((uint16(baseColor.B) * uint16(i)) / 255)
		lut[i] = image.NewUniform(color.RGBA{r, g, b, baseColor.A})
	}
	return lut
}

package render

import (
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

func loadFontFace(fontSize, dpi float64) (font.Face, error) {
	tt, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return nil, fmt.Errorf("error parsing font: %w", err)
	}
	return opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

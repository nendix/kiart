package render

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/nendix/kiart/internal/config"
)

func createDummyImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(img, image.Rect(0, 0, 50, 100), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(50, 0, 100, 100), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	return img
}

func TestRenderer_Render(t *testing.T) {
	cfg := config.AppConfig{
		Width:            50,
		FontSize:         8.0,
		DPI:              72.0,
		SkipHex:          "#000000",
		TolerancePercent: 2.0,
		OutputPath:       "unused.png",
		FontHex:          "#FFFFFF",
		BgHex:            "#000000",
	}

	renderer, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer renderer.Close()

	result := renderer.Render(createDummyImage())

	if result.Bounds().Dx() == 0 || result.Bounds().Dy() == 0 {
		t.Error("Rendered image has zero dimensions")
	}
}

func TestRenderer_RenderColored(t *testing.T) {
	cfg := config.AppConfig{
		Width:    30,
		FontSize: 8.0,
		DPI:      72.0,
		Colored:  true,
		FontHex:  "#FFFFFF",
		BgHex:    "#000000",
	}

	renderer, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer renderer.Close()

	result := renderer.Render(createDummyImage())

	if result.Bounds().Dx() == 0 || result.Bounds().Dy() == 0 {
		t.Error("Rendered image has zero dimensions")
	}
}

func TestRenderer_RenderShaded(t *testing.T) {
	cfg := config.AppConfig{
		Width:    30,
		FontSize: 8.0,
		DPI:      72.0,
		Shaded:   true,
		FontHex:  "#FF0000",
		BgHex:    "#000000",
	}

	renderer, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer renderer.Close()

	result := renderer.Render(createDummyImage())

	if result.Bounds().Dx() == 0 || result.Bounds().Dy() == 0 {
		t.Error("Rendered image has zero dimensions")
	}
}

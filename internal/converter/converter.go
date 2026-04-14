package converter

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"

	"github.com/nendix/kiart/internal/config"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var asciiChars = []rune{
	' ', '.', '\'', '`', '^', '"', ',', ':', ';', 'I', 'l', '!', 'i', '>', '<', '~', '+', '_', '-', '?',
	']', '[', '}', '{', '1', ')', '(', '|', '\\', '/', 't', 'f', 'j', 'r', 'x', 'n', 'u', 'v', 'c', 'z',
	'X', 'Y', 'U', 'J', 'C', 'L', 'Q', '0', 'O', 'Z', 'm', 'w', 'q', 'p', 'd', 'b', 'k', 'h', 'a', 'o',
	'*', '#', 'M', 'W', '&', '8', '%', 'B', '@', '$',
}

var asciiLookup [256]rune

func init() {
	for i := 0; i <= 255; i++ {
		idx := int((float64(i) / 255.0) * float64(len(asciiChars)-1))
		asciiLookup[i] = asciiChars[idx]
	}
}

// ProcessAndSave orchestrates the entire image to ASCII conversion pipeline
func ProcessAndSave(img image.Image, cfg config.AppConfig) error {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 1. Setup Skip Color Filtering
	var skipColor color.RGBA
	var actualTolerance float64
	useColorSkip := cfg.SkipHex != ""

	if useColorSkip {
		parsedColor, err := ParseHexColor(cfg.SkipHex)
		if err != nil {
			return fmt.Errorf("invalid skip hex color '%s': %w", cfg.SkipHex, err)
		}
		skipColor = parsedColor
		maxDist := math.Sqrt(255.0*255.0 + 255.0*255.0 + 255.0*255.0)
		actualTolerance = (cfg.TolerancePercent / 100.0) * maxDist
	}

	// 2. Setup Font & Dimensions
	face, err := loadFontFace(cfg.FontSize, cfg.DPI)
	if err != nil {
		return err
	}
	defer face.Close()

	metrics := face.Metrics()
	charHeight := metrics.Height.Ceil()
	advance, _ := face.GlyphAdvance('@')
	charWidth := advance.Ceil()

	newHeight := int(float64(cfg.Width) * (float64(height) / float64(width)) * (float64(charWidth) / float64(charHeight)))
	outWidth := cfg.Width * charWidth
	outHeight := newHeight * charHeight

	// 3. Setup Canvas Background
	bgColor, err := parseColorOptional(cfg.BgHex, color.RGBA{0, 0, 0, 255})
	if err != nil {
		return fmt.Errorf("invalid background hex: %w", err)
	}
	outImg := image.NewRGBA(image.Rect(0, 0, outWidth, outHeight))
	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 4. Setup Brushes
	fontColor, err := parseColorOptional(cfg.FontHex, color.RGBA{255, 255, 255, 255})
	if err != nil {
		return fmt.Errorf("invalid font hex: %w", err)
	}

	staticBrush := image.NewUniform(fontColor)
	dynamicBrush := &image.Uniform{}
	shadedUniforms := buildShadedLUT(fontColor)

	d := &font.Drawer{
		Dst:  outImg,
		Src:  staticBrush,
		Face: face,
	}

	// 5. Core Rendering Loop
	for y := range newHeight {
		for x := range cfg.Width {
			srcX := int(float64(x)/float64(cfg.Width)*float64(width)) + bounds.Min.X
			srcY := int(float64(y)/float64(newHeight)*float64(height)) + bounds.Min.Y

			pixel := img.At(srcX, srcY)

			if useColorSkip {
				pr, pg, pb, _ := pixel.RGBA()
				currentRGB := color.RGBA{uint8(pr >> 8), uint8(pg >> 8), uint8(pb >> 8), 255}
				if colorDistance(currentRGB, skipColor) <= actualTolerance {
					continue
				}
			}

			// Cleanly extract just the Y (luminance) value as a uint8
			luminance := color.GrayModel.Convert(pixel).(color.Gray).Y
			char := asciiLookup[luminance]

			if cfg.Colored {
				dynamicBrush.C = pixel
				d.Src = dynamicBrush
			} else if cfg.Shaded {
				d.Src = shadedUniforms[luminance]
			} else {
				d.Src = staticBrush
			}

			d.Dot = fixed.Point26_6{
				X: fixed.I(x * charWidth),
				Y: fixed.I((y * charHeight) + metrics.Ascent.Ceil()),
			}
			d.DrawString(string(char))
		}
	}

	// 6. Output to File
	return savePNG(outImg, cfg.OutputPath)
}

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

func parseColorOptional(hexStr string, defaultColor color.RGBA) (color.RGBA, error) {
	if hexStr == "" {
		return defaultColor, nil
	}
	return ParseHexColor(hexStr)
}

func buildShadedLUT(baseColor color.RGBA) [256]*image.Uniform {
	var lut [256]*image.Uniform
	for i := 0; i <= 255; i++ {
		r := uint8((uint16(baseColor.R) * uint16(i)) / 255)
		g := uint8((uint16(baseColor.G) * uint16(i)) / 255)
		b := uint8((uint16(baseColor.B) * uint16(i)) / 255)
		lut[i] = image.NewUniform(color.RGBA{r, g, b, baseColor.A})
	}
	return lut
}

func savePNG(img image.Image, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("error encoding PNG: %w", err)
	}
	return nil
}

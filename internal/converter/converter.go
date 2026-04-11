package converter

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Config holds all the settings for the ASCII conversion
type Config struct {
	Width            int
	FontSize         float64
	DPI              float64
	SkipHex          string
	TolerancePercent float64
	OutputPath       string
}

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

// ProcessAndSave converts the image based on the provided Config
// We return an 'error' instead of just printing to stdout so the CLI can handle it
func ProcessAndSave(img image.Image, cfg Config) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var skipColor color.RGBA
	useColorSkip := false
	var actualTolerance float64

	if cfg.SkipHex != "" {
		parsedColor, err := ParseHexColor(cfg.SkipHex)
		if err == nil {
			skipColor = parsedColor
			useColorSkip = true
			maxDist := math.Sqrt(255.0*255.0 + 255.0*255.0 + 255.0*255.0)
			actualTolerance = (cfg.TolerancePercent / 100.0) * maxDist
		} else {
			// In a real app, you might want to return this error, but we'll just warn
			fmt.Printf("Warning: Could not parse hex color '%s'. Color skipping disabled.\n", cfg.SkipHex)
		}
	}

	tt, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return fmt.Errorf("error parsing font: %w", err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    cfg.FontSize,
		DPI:     cfg.DPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("error creating font face: %w", err)
	}
	defer face.Close()

	metrics := face.Metrics()
	charHeight := metrics.Height.Ceil()
	advance, _ := face.GlyphAdvance('@')
	charWidth := advance.Ceil()

	fontRatio := float64(charWidth) / float64(charHeight)
	imgRatio := float64(height) / float64(width)
	newHeight := int(float64(cfg.Width) * imgRatio * fontRatio)

	outWidth := cfg.Width * charWidth
	outHeight := newHeight * charHeight
	outImg := image.NewRGBA(image.Rect(0, 0, outWidth, outHeight))
	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  outImg,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	for y := range newHeight {
		for x := range cfg.Width {
			srcX := int(float64(x)/float64(cfg.Width)*float64(width)) + bounds.Min.X
			srcY := int(float64(y)/float64(newHeight)*float64(height)) + bounds.Min.Y

			pixel := img.At(srcX, srcY)

			if useColorSkip {
				pr, pg, pb, _ := pixel.RGBA()
				currentRGB := color.RGBA{uint8(pr >> 8), uint8(pg >> 8), uint8(pb >> 8), 255}

				dist := colorDistance(currentRGB, skipColor)

				if dist <= actualTolerance {
					continue
				}
			}

			grayColor := color.GrayModel.Convert(pixel).(color.Gray)
			char := asciiLookup[grayColor.Y]

			d.Dot = fixed.Point26_6{
				X: fixed.I(x * charWidth),
				Y: fixed.I((y * charHeight) + metrics.Ascent.Ceil()),
			}
			d.DrawString(string(char))
		}
	}

	outFile, err := os.Create(cfg.OutputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, outImg)
	if err != nil {
		return fmt.Errorf("error encoding PNG: %w", err)
	}

	fmt.Printf("Successfully generated ASCII art (%dx%d chars) -> %s\n", cfg.Width, newHeight, cfg.OutputPath)
	return nil
}

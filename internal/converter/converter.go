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
	Shaded           bool
	Colored          bool
	FontHex          string
	BgHex            string
}

func DefaultConfig() Config {
	return Config{
		Width:            510,
		FontSize:         8.0,
		DPI:              300.0,
		SkipHex:          "",
		TolerancePercent: 2.0,
		OutputPath:       "ascii_art.png",
		Shaded:           false,
		Colored:          false,
		FontHex:          "#FFFFFF",
		BgHex:            "#000000",
	}
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

func ProcessAndSave(img image.Image, cfg Config) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var skipColor color.RGBA
	useColorSkip := false
	var actualTolerance float64

	if cfg.SkipHex != "" {
		parsedColor, err := ParseHexColor(cfg.SkipHex)
		if err != nil {
			return fmt.Errorf("invalid skip hex color '%s': %w", cfg.SkipHex, err)
		}

		skipColor = parsedColor
		useColorSkip = true
		maxDist := math.Sqrt(255.0*255.0 + 255.0*255.0 + 255.0*255.0)
		actualTolerance = (cfg.TolerancePercent / 100.0) * maxDist
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

	bgColor := color.RGBA{0, 0, 0, 255} // Default to solid black
	if cfg.BgHex != "" {
		parsedBg, err := ParseHexColor(cfg.BgHex)
		if err != nil {
			return fmt.Errorf("invalid background hex color '%s': %w", cfg.BgHex, err)
		}
		bgColor = parsedBg
	}

	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	brushColor := color.RGBA{255, 255, 255, 255}
	if cfg.FontHex != "" {
		parsedFontColor, err := ParseHexColor(cfg.FontHex)
		if err != nil {
			return fmt.Errorf("invalid font hex color '%s': %w", cfg.FontHex, err)
		}
		brushColor = parsedFontColor
	}
	staticBrush := image.NewUniform(brushColor)
	dynamicBrush := &image.Uniform{}

	var shadedUniforms [256]*image.Uniform
	for i := 0; i <= 255; i++ {
		r := uint8((uint16(brushColor.R) * uint16(i)) / 255)
		g := uint8((uint16(brushColor.G) * uint16(i)) / 255)
		b := uint8((uint16(brushColor.B) * uint16(i)) / 255)
		a := brushColor.A

		shadedUniforms[i] = image.NewUniform(color.RGBA{r, g, b, a})
	}

	d := &font.Drawer{
		Dst:  outImg,
		Src:  staticBrush,
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

			luminance := color.GrayModel.Convert(pixel).(color.Gray)
			char := asciiLookup[luminance.Y]

			if cfg.Colored {
				dynamicBrush.C = pixel
				d.Src = dynamicBrush
			} else if cfg.Shaded {
				d.Src = shadedUniforms[luminance.Y]
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

	outFile, err := os.Create(cfg.OutputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, outImg)
	if err != nil {
		return fmt.Errorf("error encoding PNG: %w", err)
	}

	// It's still okay to print success here, but often in Go libraries,
	// even the success print is pushed to main.go. I've left it as is for your CLI experience!
	fmt.Printf("Successfully generated ASCII art (%dx%d chars) -> %s\n", cfg.Width, newHeight, cfg.OutputPath)
	return nil
}

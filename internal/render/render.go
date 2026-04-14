package render

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/nendix/kiart/internal/ascii"
	kicolor "github.com/nendix/kiart/internal/color"
	"github.com/nendix/kiart/internal/config"
)

// Renderer converts images to ASCII art.
type Renderer struct {
	cfg             config.AppConfig
	face            font.Face
	charWidth       int
	charHeight      int
	ascent          int
	fontColor       color.RGBA
	bgColor         color.RGBA
	skipColor       color.RGBA
	actualTolerance float64
	useColorSkip    bool
}

// New creates a Renderer, validating config and loading font resources.
func New(cfg config.AppConfig) (*Renderer, error) {
	face, err := loadFontFace(cfg.FontSize, cfg.DPI)
	if err != nil {
		return nil, err
	}

	metrics := face.Metrics()
	charHeight := metrics.Height.Ceil()
	advance, _ := face.GlyphAdvance('@')
	charWidth := advance.Ceil()

	fontColor, err := kicolor.ParseHexWithDefault(cfg.FontHex, color.RGBA{255, 255, 255, 255})
	if err != nil {
		return nil, fmt.Errorf("invalid font hex: %w", err)
	}

	bgColor, err := kicolor.ParseHexWithDefault(cfg.BgHex, color.RGBA{0, 0, 0, 255})
	if err != nil {
		return nil, fmt.Errorf("invalid background hex: %w", err)
	}

	r := &Renderer{
		cfg:        cfg,
		face:       face,
		charWidth:  charWidth,
		charHeight: charHeight,
		ascent:     metrics.Ascent.Ceil(),
		fontColor:  fontColor,
		bgColor:    bgColor,
	}

	if cfg.SkipHex != "" {
		skipColor, err := kicolor.ParseHex(cfg.SkipHex)
		if err != nil {
			return nil, fmt.Errorf("invalid skip hex color '%s': %w", cfg.SkipHex, err)
		}
		r.skipColor = skipColor
		r.useColorSkip = true
		maxDist := math.Sqrt(255.0*255.0 + 255.0*255.0 + 255.0*255.0)
		r.actualTolerance = (cfg.TolerancePercent / 100.0) * maxDist
	}

	return r, nil
}

// Render converts an input image to ASCII art, returning the rendered image.
func (r *Renderer) Render(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	srcWidth, srcHeight := bounds.Dx(), bounds.Dy()

	newHeight := int(float64(r.cfg.Width) * (float64(srcHeight) / float64(srcWidth)) * (float64(r.charWidth) / float64(r.charHeight)))
	outWidth := r.cfg.Width * r.charWidth
	outHeight := newHeight * r.charHeight

	outImg := image.NewRGBA(image.Rect(0, 0, outWidth, outHeight))
	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{r.bgColor}, image.Point{}, draw.Src)

	staticBrush := image.NewUniform(r.fontColor)
	dynamicBrush := &image.Uniform{}
	shadedLUT := buildShadedLUT(r.fontColor)

	d := &font.Drawer{
		Dst:  outImg,
		Src:  staticBrush,
		Face: r.face,
	}

	for y := range newHeight {
		for x := range r.cfg.Width {
			srcX := int(float64(x)/float64(r.cfg.Width)*float64(srcWidth)) + bounds.Min.X
			srcY := int(float64(y)/float64(newHeight)*float64(srcHeight)) + bounds.Min.Y

			pixel := img.At(srcX, srcY)

			if r.useColorSkip {
				pr, pg, pb, _ := pixel.RGBA()
				currentRGB := color.RGBA{uint8(pr >> 8), uint8(pg >> 8), uint8(pb >> 8), 255}
				if kicolor.Distance(currentRGB, r.skipColor) <= r.actualTolerance {
					continue
				}
			}

			luminance := color.GrayModel.Convert(pixel).(color.Gray).Y
			char := ascii.Lookup[luminance]

			if r.cfg.Colored {
				dynamicBrush.C = pixel
				d.Src = dynamicBrush
			} else if r.cfg.Shaded {
				d.Src = shadedLUT[luminance]
			} else {
				d.Src = staticBrush
			}

			d.Dot = fixed.Point26_6{
				X: fixed.I(x * r.charWidth),
				Y: fixed.I(y*r.charHeight + r.ascent),
			}
			d.DrawString(string(char))
		}
	}

	return outImg
}

// RenderText converts an input image to plain ASCII text.
func (r *Renderer) RenderText(img image.Image) string {
	bounds := img.Bounds()
	srcWidth, srcHeight := bounds.Dx(), bounds.Dy()

	// Terminal chars are ~2x taller than wide, so halve the height ratio.
	newHeight := int(float64(r.cfg.Width) * (float64(srcHeight) / float64(srcWidth)) * 0.5)

	var buf strings.Builder
	for y := range newHeight {
		for x := range r.cfg.Width {
			srcX := int(float64(x)/float64(r.cfg.Width)*float64(srcWidth)) + bounds.Min.X
			srcY := int(float64(y)/float64(newHeight)*float64(srcHeight)) + bounds.Min.Y

			pixel := img.At(srcX, srcY)

			if r.useColorSkip {
				pr, pg, pb, _ := pixel.RGBA()
				currentRGB := color.RGBA{uint8(pr >> 8), uint8(pg >> 8), uint8(pb >> 8), 255}
				if kicolor.Distance(currentRGB, r.skipColor) <= r.actualTolerance {
					buf.WriteByte(' ')
					continue
				}
			}

			luminance := color.GrayModel.Convert(pixel).(color.Gray).Y
			buf.WriteRune(ascii.Lookup[luminance])
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

// Close releases font resources.
func (r *Renderer) Close() error {
	return r.face.Close()
}

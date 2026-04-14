package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/nendix/kiart/internal/config"
	"github.com/nendix/kiart/internal/converter"
)

func main() {
	cfg := config.NewDefault()

	flag.IntVarP(&cfg.Width, "width", "w", cfg.Width, "Width of the output ASCII art in characters")
	flag.Float64VarP(&cfg.FontSize, "font-size", "s", cfg.FontSize, "Font size in points")
	flag.Float64VarP(&cfg.DPI, "dpi", "d", cfg.DPI, "Resolution in Dots Per Inch")

	flag.StringVarP(&cfg.SkipHex, "skip", "", cfg.SkipHex, "HEX color to filter out")
	flag.Float64VarP(&cfg.TolerancePercent, "tolerance", "t", cfg.TolerancePercent, "Color tolerance percentage (0-100)")

	flag.StringVarP(&cfg.OutputPath, "output", "o", cfg.OutputPath, "Path to save the output image")

	flag.BoolVarP(&cfg.Shaded, "shaded", "S", cfg.Shaded, "Render characters using true grayscale shading")
	flag.BoolVarP(&cfg.Colored, "colored", "C", cfg.Colored, "Render characters using original RGB colors")

	flag.StringVarP(&cfg.FontHex, "font-color", "c", cfg.FontHex, "Font HEX color")
	flag.StringVarP(&cfg.BgHex, "background", "b", cfg.BgHex, "Canvas background HEX color or 'transparent'")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "kiart - Image to ASCII Art Converter\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  kiart [options] <path-to-image>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")

		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	imagePath := args[0]

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening image: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding image: %v\n", err)
		os.Exit(1)
	}

	err = converter.ProcessAndSave(img, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated ASCII art -> %s\n", cfg.OutputPath)
}

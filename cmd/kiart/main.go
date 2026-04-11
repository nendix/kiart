package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/nendix/kiart/internal/converter"
)

func main() {
	var cfg converter.Config

	// Syntax: (&variable, "long-name", "short-name", default_value, "Description")
	flag.IntVarP(&cfg.Width, "width", "w", 510, "Width of the output ASCII art in characters")
	flag.Float64VarP(&cfg.FontSize, "size", "s", 8.0, "Font size in points")
	flag.Float64VarP(&cfg.DPI, "dpi", "d", 300.0, "Resolution in Dots Per Inch")

	flag.StringVarP(&cfg.SkipHex, "skip", "", "", "HEX color to filter out and make transparent")
	flag.Float64VarP(&cfg.TolerancePercent, "tol", "", 2.0, "Color tolerance percentage (0-100)")

	flag.StringVarP(&cfg.OutputPath, "out", "o", "ascii_art.png", "Path to save the output image")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "kiart - Blazing fast, high-res ASCII art generator\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  kiart [options] <path-to-image>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")

		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  kiart image.jpg\n")
		fmt.Fprintf(os.Stderr, "  kiart -w 120 -s 14 image.jpg\n")
		fmt.Fprintf(os.Stderr, "  kiart --skip #3A3A3A --tol 5 -o ascii_art.png image.jpg\n")
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
		fmt.Printf("Error opening image: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		os.Exit(1)
	}

	err = converter.ProcessAndSave(img, cfg)
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}

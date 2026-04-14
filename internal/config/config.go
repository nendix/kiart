package config

type AppConfig struct {
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

func NewDefault() AppConfig {
	return AppConfig{
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

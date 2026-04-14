# kiart

A fast CLI tool to convert images into ASCII art. Outputs plain text to your terminal or renders high-resolution PNG images.

## Install

**Homebrew:**

```bash
brew tap nendix/tap
brew install kiart
```

**Go:**

```bash
go install github.com/nendix/kiart/cmd/kiart@latest
```

**From source:**

```bash
git clone https://github.com/nendix/kiart.git
cd kiart
make install
```

## Usage

```bash
# Print ASCII art to terminal
kiart photo.jpg

# Save as PNG image
kiart -o output.png photo.jpg

# Colored output preserving original colors
kiart -C -o output.png photo.jpg

# Grayscale shading with custom font color
kiart -S -c "#00FF00" -o output.png photo.jpg

# Transparent background, custom width
kiart -w 200 -b transparent -o output.png photo.jpg

# Skip a background color (e.g. remove white)
kiart --skip "#FFFFFF" -t 5 -o output.png photo.jpg
```

## Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--width` | `-w` | `480` | Width of output in characters |
| `--output` | `-o` | stdout | Path to save output PNG |
| `--colored` | `-C` | `false` | Render with original RGB colors |
| `--shaded` | `-S` | `false` | Render with grayscale shading |
| `--font-color` | `-c` | `#FFFFFF` | Font color (hex) |
| `--background` | `-b` | `#000000` | Background color (hex or `transparent`) |
| `--font-size` | `-s` | `8.0` | Font size in points |
| `--dpi` | `-d` | `150` | Resolution in DPI |
| `--skip` | | | Hex color to filter out |
| `--tolerance` | `-t` | `2.0` | Color skip tolerance (0-100%) |
| `--version` | `-v` | | Print version and exit |

## Rendering Modes

- **Default** -- White characters on black background. Each character chosen by luminance.
- **Shaded** (`-S`) -- Characters tinted by luminance. Brighter areas get brighter characters.
- **Colored** (`-C`) -- Each character rendered in the original pixel color.

When no `-o` flag is given, kiart prints plain ASCII text to stdout. With `-o`, it renders a high-resolution PNG using the embedded GoMono font.

## Supported Formats

Input: JPEG, PNG

Output: PNG (with `-o`) or plain text (stdout)

## License

MIT

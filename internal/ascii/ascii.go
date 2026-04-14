package ascii

// Charset ordered by visual density from lightest to heaviest.
var Charset = []rune{
	' ', '.', '\'', '`', '^', '"', ',', ':', ';', 'I', 'l', '!', 'i', '>', '<', '~', '+', '_', '-', '?',
	']', '[', '}', '{', '1', ')', '(', '|', '\\', '/', 't', 'f', 'j', 'r', 'x', 'n', 'u', 'v', 'c', 'z',
	'X', 'Y', 'U', 'J', 'C', 'L', 'Q', '0', 'O', 'Z', 'm', 'w', 'q', 'p', 'd', 'b', 'k', 'h', 'a', 'o',
	'*', '#', 'M', 'W', '&', '8', '%', 'B', '@', '$',
}

// Lookup maps a luminance value (0-255) to an ASCII character.
var Lookup [256]rune

func init() {
	for i := range 256 {
		idx := int((float64(i) / 255.0) * float64(len(Charset)-1))
		Lookup[i] = Charset[idx]
	}
}

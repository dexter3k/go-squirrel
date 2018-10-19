package lexer

func isExponent(r rune) bool {
	return r == 'e' || r == 'E'
}

func isOctal(r rune) bool {
	return r >= '0' && r <= '7'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHex(r rune) bool {
	return isDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isAlnum(r rune) bool {
	return isDigit(r) || isAlpha(r)
}

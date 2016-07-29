package hello8

// %b: -ddddddddp±ddd
func fmtB(dst []byte, neg bool, mant uint64, exp int, flt *floatInfo) []byte {
	// sign
	if neg {
		dst = append(dst, '-')
	}

	// mantissa
	dst, _ = formatBits(dst, mant, 10, false, true)

	// p
	dst = append(dst, 'p')

	// ±exponent
	exp -= int(flt.mantbits)
	if exp >= 0 {
		dst = append(dst, '+')
	}
	dst, _ = formatBits(dst, uint64(exp), 10, exp < 0, true)

	return dst
}

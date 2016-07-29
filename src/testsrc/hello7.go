package hello7

// Frexp10 is an analogue of math.Frexp for decimal powers. It scales
// f by an approximate power of ten 10^-exp, and returns exp10, so
// that f*10^exp10 has the same value as the old f, up to an ulp,
// as well as the index of 10^-exp in the powersOfTen table.
func (f *extFloat) frexp10() (exp10, index int) {
	// The constants expMin and expMax constrain the final value of the
	// binary exponent of f. We want a small integral part in the result
	// because finding digits of an integer requires divisions, whereas
	// digits of the fractional part can be found by repeatedly multiplying
	// by 10.
	const expMin = -60
	const expMax = -32
	// Find power of ten such that x * 10^n has a binary exponent
	// between expMin and expMax.
	approxExp10 := ((expMin+expMax)/2 - f.exp) * 28 / 93 // log(10)/log(2) is close to 93/28.
	i := (approxExp10 - firstPowerOfTen) / stepPowerOfTen
Loop:
	for {
		exp := f.exp + powersOfTen[i].exp + 64
		switch {
		case exp < expMin:
			i++
		case exp > expMax:
			i--
		default:
			break Loop
		}
	}
	// Apply the desired decimal shift on f. It will have exponent
	// in the desired range. This is multiplication by 10^-exp10.
	f.Multiply(powersOfTen[i])

	return -(firstPowerOfTen + i*stepPowerOfTen), i
}

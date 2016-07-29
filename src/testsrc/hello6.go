package hello6

// RuneCount returns the number of runes in p. Erroneous and short
// encodings are treated as single runes of width 1 byte.
func RuneCount(p []byte) int {
	np := len(p)
	var n int
	for i := 0; i < np; {
		n++
		c := p[i]
		if c < RuneSelf {
			// ASCII fast path
			i++
			continue
		}
		x := first[c]
		if x == xx {
			i++ // invalid.
			continue
		}
		size := int(x & 7)
		if i+size > np {
			i++ // Short or invalid.
			continue
		}
		accept := acceptRanges[x>>4]
		if c := p[i+1]; c < accept.lo || accept.hi < c {
			size = 1
		} else if size == 2 {
		} else if c := p[i+2]; c < locb || hicb < c {
			size = 1
		} else if size == 3 {
		} else if c := p[i+3]; c < locb || hicb < c {
			size = 1
		}
		i += size
	}
	return n
}

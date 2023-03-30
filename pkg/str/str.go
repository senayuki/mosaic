package str

import (
	"unicode/utf8"
)

func Runes2Bytes(in []rune) []byte {
	out := make([]byte, len(in)*utf8.UTFMax)
	count := 0
	for _, r := range in {
		count += utf8.EncodeRune(out[count:], r)
	}
	out = out[:count]
	return out
}

func Bytes2Runes(in []byte) []rune {
	out := make([]rune, 0, len(in))
	ptr := 0
	for {
		r, size := utf8.DecodeRune(in[ptr:])
		out = append(out, r)
		ptr += size
		if len(in) <= ptr {
			break
		}
	}
	return out
}

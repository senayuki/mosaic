package mask

import "github.com/senayuki/mosaic/types"

type MarkCoverProcesser struct {
	CoverChar rune
	Offset    int
	Padding   int
	Length    int
	Reverse   bool
}

func (m MarkCoverProcesser) DefaultCoverChar() rune {
	return rune('*')
}

func (m *MarkCoverProcesser) Init(maskRule *types.KVMaskConfig) {
	coverParam := maskRule.CoverParam
	m.CoverChar = m.DefaultCoverChar()
	if len(coverParam.Char) > 0 {
		m.CoverChar = []rune(coverParam.Char)[0]
	}
	m.Offset = maskRule.CoverParam.Offset
	m.Padding = maskRule.CoverParam.Padding
	m.Length = maskRule.CoverParam.Length
	m.Reverse = maskRule.CoverParam.Reverse
}

// TODO input should be []byte
func (m MarkCoverProcesser) Mask(in string) (out string, err error) {
	inRune := []rune(in)
	inLen := len(inRune)
	if len(in) == 0 {
		return
	}
	outRune := []rune{}
	maskedLength := 0
	offset := m.Offset
	padding := m.Padding
	if inLen <= offset && inLen <= padding {
		offset = (inLen - 1) / 2
		padding = (inLen - 1) / 2
	} else if inLen <= offset {
		offset = inLen - 1
	} else if inLen <= padding {
		padding = inLen - 1
	}
	for index, _ := range inRune {
		if index >= offset && index < (inLen-padding) {
			if m.Length > 0 {
				if maskedLength < m.Length {
					outRune = append(outRune, m.CoverChar)
					maskedLength++
				} else {
					// NOOP, skip the rune
				}
			} else {
				outRune = append(outRune, m.CoverChar)
				maskedLength++
			}
		} else {
			outRune = append(outRune, inRune[index])
		}
	}
	return string(outRune), nil
}

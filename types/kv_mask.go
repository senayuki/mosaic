package types

type KVMaskConfig struct {
	RuleName   string
	MaskType   MaskType
	CoverParam MaskRuleCoverParam
}

type (
	MaskType string
)

const (
	MaskTypeCover MaskType = "cover"
)

type MaskRuleCoverParam struct {
	/*cover by char，*，0，# etc.
	* is default value
	if length of string more than 1, use the first char
	*/
	Char string
	/*cover start at this index
	start at 0 by default
	if offset is bigger than length of content, will ensure that last char of content covered at least
	if offset & padding is overlapped, all chars will be covered
	*/
	Offset int
	/*cover offset from tail
	end at last char by default
	if padding is bigger than length of content, will ensure that first char of content covered at least
	if offset & padding is overlapped, all chars will be covered
	*/
	Padding int
	/*length of covered string
	if value equals 0, the length decide by content
	if value more than length that can be mask, sames like equals 0
	*/
	Length int
	/*cover string start at the tail*/
	Reverse bool
}

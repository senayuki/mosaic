package types

type KVMaskConfig struct {
	RuleName string
	MaskType MaskType
	Params   interface{}
}

type (
	MaskType string
)

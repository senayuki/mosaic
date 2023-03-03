package types

type Mask struct {
	JSONPath   []string
	Text       interface{}
	MaskedText interface{}
}

type DLPRule struct {
	Detect struct{}
}

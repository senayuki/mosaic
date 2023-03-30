package types

type (
	KVDetectConfig struct {
		KeyEqs      []string    // key absolutely equal an element in array
		ValEqs      []string    // val absolutely equal an element in array
		KeyContains []string    // key contains an element in array
		ValContains []string    // val contains an element in array
		KeyRegex    []string    // keys matched an regex
		ValRegex    []string    // vals matched an regex
		MatchMode   KVMatchMode // (key || val) matched or (key && val) matched
		ValueMode   KVMaskMode  // mask whole value or matched segments
		MaskRef     string      // reference mask processes
		/*treat specified field as key-value pair
		{
			"name": "as key, val must be string",
			"content": "as value, value can be any type exclude object"
		} */
		KVFieldOpt *KVField
	}
	KVMatchMode string // (key || val) matched or (key && val) matched
	KVMaskMode  string // mask whole value or matched value
	KVField     struct {
		Key string
		Val string
	}
)

const (
	KVMatchDefault KVMatchMode = ""    // "or" is default mode
	KVMatchOr      KVMatchMode = "or"  // key or val matched
	KVMatchAnd     KVMatchMode = "and" // key and val matched

	// TODO: mask mode support
	KVMaskModeDefault KVMaskMode = ""        // "whole" is default mode
	KVMaskModeWhole   KVMaskMode = "whole"   // whole value
	KVMaskModeSegment KVMaskMode = "segment" // matched segments in value
)

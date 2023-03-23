package types

import "fmt"

type KVPair struct {
	Key         string
	Val         interface{}
	ValJSONPath JSONPath
	ValMasked   interface{}
	KVFieldRel  *KVField
}

func (kv *KVPair) GetValString() string {
	if kv.Val == nil {
		return "null"
	}
	return fmt.Sprintf("%v", kv.Val)
}

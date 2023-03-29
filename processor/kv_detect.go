package processer

import (
	"regexp"
	"strings"

	"github.com/senayuki/mosaic/types"
	"github.com/valyala/fastjson"
)

type detectExp struct {
	KeyRegex []*regexp.Regexp
	ValRegex []*regexp.Regexp
}

// input JSON bytes
func (m KVProcesser) Detect(input []byte) ([]types.KVPair, error) {
	val, err := fastjson.ParseBytes(input)
	if err != nil {
		return nil, err
	}
	matched := []types.KVPair{}
	// extract all k-v pair (include k-v pair in fields
	elements := m.visit(types.NewJSONPath(), "", val, nil)
	for _, v := range elements {
		valString := v.GetValString()
		for configIdx, config := range m.detectConfig {
			if v.KVFieldRel != config.KVFieldOpt {
				continue
			}
			if m.matchKV(configIdx, v, valString) {
				matched = append(matched, v)
			}
		}
	}
	return matched, nil
}

func (m KVProcesser) matchKV(configIdx int, pair types.KVPair, valString string) bool {
	keyEqMatch := false
	keyContainsMatch := false
	keyRegMatch := false
	for _, keyKeyword := range m.detectConfig[configIdx].KeyEqs {
		if strings.EqualFold(keyKeyword, pair.Key) {
			keyEqMatch = true
			break
		}
	}
	for _, keyContains := range m.detectConfig[configIdx].KeyContains {
		if strings.Contains(keyContains, pair.Key) {
			keyContainsMatch = true
			break
		}
	}
	for _, keyRegex := range m.detectExp[configIdx].KeyRegex {
		if keyRegex.MatchString(pair.Key) {
			keyRegMatch = true
			break
		}
	}

	valEqMatch := false
	valContainsMatch := false
	valRegMatch := false
	for _, valKeyword := range m.detectConfig[configIdx].ValEqs {
		if strings.EqualFold(valKeyword, valString) {
			valEqMatch = true
			break
		}
	}
	for _, valContains := range m.detectConfig[configIdx].ValContains {
		if strings.Contains(valContains, valString) {
			valContainsMatch = true
			break
		}
	}
	for _, valRegex := range m.detectExp[configIdx].ValRegex {
		if valRegex.MatchString(valString) {
			valRegMatch = true
			break
		}
	}
	switch m.detectConfig[configIdx].MatchMode {
	case types.KVMatchDefault, types.KVMatchOr:
		return (keyEqMatch || keyContainsMatch || keyRegMatch) ||
			(valEqMatch || valContainsMatch || valRegMatch)
	case types.KVMatchAnd:
		return (keyEqMatch || keyContainsMatch || keyRegMatch) &&
			(valEqMatch || valContainsMatch || valRegMatch)
	default:
		return false
	}
}

// recursion to extract elements
func (m KVProcesser) visit(valJSONPath types.JSONPath, key string, val *fastjson.Value, kvFieldRel *types.KVField) []types.KVPair {
	elements := []types.KVPair{}
	switch val.Type() {
	case fastjson.TypeObject:
		obj := val.GetObject()
		keyFieldProbable := map[string]struct{}{}
		field := map[string]*fastjson.Value{}
		obj.Visit(func(key []byte, v *fastjson.Value) {
			keyStr := string(key)
			if _, ok := m.detectKVField[keyStr]; ok {
				keyFieldProbable[keyStr] = struct{}{}
			}
			field[keyStr] = v
			elements = append(elements, m.visit(valJSONPath.Append(keyStr), keyStr, v, nil)...)
		})
		// add k-v fields relations
		if len(keyFieldProbable) > 0 {
			for keyField, _ := range keyFieldProbable {
				// key must be string
				if keyVal, ok := field[keyField]; ok && keyVal.Type() == fastjson.TypeString {
					// value of keyField is the real key
					key := string(keyVal.GetStringBytes())
					for valField, kvFieldRel := range m.detectKVField[keyField] {
						// value is impossible to be object
						if val, ok := field[valField]; ok && val.Type() != fastjson.TypeObject {
							elements = append(elements, m.visit(valJSONPath.Append(valField), key, val, kvFieldRel)...)
						}
					}
				}
			}
		}
	case fastjson.TypeArray:
		arr := val.GetArray()
		for idx, item := range arr {
			elements = append(elements, m.visit(valJSONPath.Append(idx), key, item, kvFieldRel)...)
		}
	case fastjson.TypeString:
		elements = append(elements, types.KVPair{
			Key:         key,
			ValJSONPath: valJSONPath,
			Val:         string(val.GetStringBytes()),
			ValMasked:   string(val.GetStringBytes()),
			KVFieldRel:  kvFieldRel,
		})
	case fastjson.TypeNumber:
		elements = append(elements, types.KVPair{
			Key:         key,
			ValJSONPath: valJSONPath,
			Val:         val.GetInt64(),
			ValMasked:   val.GetInt64(),
			KVFieldRel:  kvFieldRel,
		})
	case fastjson.TypeNull:
		elements = append(elements, types.KVPair{
			Key:         key,
			ValJSONPath: valJSONPath,
			Val:         nil,
			ValMasked:   nil,
			KVFieldRel:  kvFieldRel,
		})
	case fastjson.TypeTrue:
		elements = append(elements, types.KVPair{
			Key:         key,
			ValJSONPath: valJSONPath,
			Val:         true,
			ValMasked:   true,
			KVFieldRel:  kvFieldRel,
		})
	case fastjson.TypeFalse:
		elements = append(elements, types.KVPair{
			Key:         key,
			ValJSONPath: valJSONPath,
			Val:         false,
			ValMasked:   false,
			KVFieldRel:  kvFieldRel,
		})
	}
	return elements
}

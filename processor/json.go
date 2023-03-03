package processer

import (
	"log"
	"strconv"

	"github.com/senayuki/mosaic/types"
	"github.com/valyala/fastjson"
)

// input JSON string
func MaskJSON(rule types.DLPRule, input string) (string, error) {
	result, err := fastjson.Parse(input)
	if err != nil {
		return "", err
	}
	matched := fastjsonDetector([]string{}, result)
	log.Println(matched)

	return result.String(), nil
}

// recursion to detect all matched items
func fastjsonDetector(jsonPath []string, val *fastjson.Value) []types.Mask {
	matched := []types.Mask{}
	switch val.Type() {
	case fastjson.TypeObject:
		obj := val.GetObject()
		obj.Visit(func(key []byte, v *fastjson.Value) {
			matched = append(matched, fastjsonDetector(addJSONPath(jsonPath, string(key)), v)...)
		})
	case fastjson.TypeArray:
		arr := val.GetArray()
		for idx, item := range arr {
			matched = append(matched, fastjsonDetector(addJSONPath(jsonPath, strconv.Itoa(idx)), item)...)
		}
	case fastjson.TypeString:
		matched = append(matched, types.Mask{
			JSONPath:   jsonPath,
			Text:       string(val.GetStringBytes()),
			MaskedText: string(val.GetStringBytes()),
		})
	case fastjson.TypeNumber:
		matched = append(matched, types.Mask{
			JSONPath:   jsonPath,
			Text:       val.GetInt64(),
			MaskedText: val.GetInt64(),
		})
	case fastjson.TypeNull:
		matched = append(matched, types.Mask{
			JSONPath:   jsonPath,
			Text:       nil,
			MaskedText: nil,
		})
	case fastjson.TypeTrue:
		matched = append(matched, types.Mask{
			JSONPath:   jsonPath,
			Text:       true,
			MaskedText: true,
		})
	case fastjson.TypeFalse:
		matched = append(matched, types.Mask{
			JSONPath:   jsonPath,
			Text:       false,
			MaskedText: false,
		})
	}
	return matched
}

func addJSONPath(old []string, add string) []string {
	jsonPathNow := make([]string, len(old), len(old)+1)
	copy(jsonPathNow, old)
	jsonPathNow = append(jsonPathNow, add)
	return jsonPathNow
}

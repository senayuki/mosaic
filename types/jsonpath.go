package types

import (
	"encoding/json"
	"strconv"
	"strings"
)

type JSONPath struct {
	path []interface{}
}

func NewJSONPath() JSONPath {
	return JSONPath{path: []interface{}{}}
}
func (j JSONPath) Append(key ...interface{}) JSONPath {
	newPath := make([]interface{}, len(j.path), len(j.path)+len(key))
	copy(newPath, j.path)
	newPath = append(newPath, key...)
	return JSONPath{path: newPath}
}
func (j JSONPath) String() string {
	return strings.Join(j.ToStrings(), "->")
}

func (j JSONPath) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.path)
}

func (j *JSONPath) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &j.path)
	return err
}

func (j JSONPath) ToStrings() []string {
	result := make([]string, 0, len(j.path))
	for _, path := range j.path {
		switch val := path.(type) {
		case string:
			result = append(result, val)
		case int:
			result = append(result, strconv.Itoa(val))
		}
	}
	return result
}

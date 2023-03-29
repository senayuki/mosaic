package processer

import (
	"regexp"

	"github.com/senayuki/mosaic/types"
)

type KVProcesser struct {
	detectConfig  []types.KVDetectConfig
	detectKVField map[string]map[string]*types.KVField // key:val fields in config
	detectExp     []detectExp                          // compiled regex
}

func NewKVProcesser(rules types.KVRules) KVProcesser {
	m := KVProcesser{detectConfig: rules.DetectRules, detectKVField: map[string]map[string]*types.KVField{}}
	for idx, config := range m.detectConfig {
		// find all KVField
		if config.KVFieldOpt != nil {
			if _, ok := m.detectKVField[config.KVFieldOpt.Key]; !ok {
				m.detectKVField[config.KVFieldOpt.Key] = map[string]*types.KVField{}
			}
			m.detectKVField[config.KVFieldOpt.Key][config.KVFieldOpt.Val] = config.KVFieldOpt
		}
		// compile regexp
		m.detectExp = append(m.detectExp, detectExp{})
		for _, regex := range config.KeyRegex {
			m.detectExp[idx].KeyRegex = append(m.detectExp[idx].KeyRegex, regexp.MustCompile(regex))
		}
		for _, regex := range config.ValRegex {
			m.detectExp[idx].ValRegex = append(m.detectExp[idx].ValRegex, regexp.MustCompile(regex))
		}
	}
	return m
}

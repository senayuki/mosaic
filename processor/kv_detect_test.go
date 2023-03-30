package processer

import (
	"encoding/json"
	"testing"

	"github.com/senayuki/mosaic/types"
	"github.com/valyala/fastjson"
)

func TestKVProcesser_Detect(t *testing.T) {
	type args struct {
		rule  types.KVRules
		input interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantDetect []types.KVPair
		wantErr    bool
	}{
		{
			name: "match normal key-val(string)",
			args: args{
				rule: types.KVRules{
					DetectRules: []types.KVDetectConfig{
						{
							KeyEqs: []string{
								"password",
							},
						},
					},
				},
				input: map[string]interface{}{
					"password":    "val1",
					"otherKey":    "1234567890",
					"otherKeyInt": 1234567890,
					"obj": map[string]interface{}{
						"password":    "val2",
						"otherKey":    "1234567890",
						"otherKeyInt": 1234567890,
					},
				},
			},
			wantDetect: []types.KVPair{
				{
					Key:         "password",
					Val:         "val2",
					ValMasked:   "val2",
					ValJSONPath: types.NewJSONPath().Append("obj").Append("password"),
					KVFieldRel:  nil,
				},
				{
					Key:         "password",
					Val:         "val1",
					ValMasked:   "val1",
					ValJSONPath: types.NewJSONPath().Append("password"),
					KVFieldRel:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "match key-val(non-string)",
			args: args{
				rule: types.KVRules{
					DetectRules: []types.KVDetectConfig{
						{
							KeyContains: []string{
								"phone",
							},
						},
						{
							ValContains: []string{
								"tru",
							},
						},
						{
							ValEqs: []string{
								"null",
							},
						},
						{
							KeyRegex:  []string{"^mobile[0-9]*$"},
							ValRegex:  []string{"^[0-9]*$"},
							MatchMode: types.KVMatchAnd,
						},
					},
				},
				input: map[string]interface{}{
					"phonenumber": 1234567890,
					"mobile1234":  "12344321",
					"isMember":    true,
					"notMember":   false,
					"status":      nil,
				},
			},
			wantDetect: []types.KVPair{
				{
					Key:         "isMember",
					Val:         true,
					ValMasked:   true,
					ValJSONPath: types.NewJSONPath().Append("isMember"),
					KVFieldRel:  nil,
				},
				{
					Key:         "mobile1234",
					Val:         "12344321",
					ValJSONPath: types.NewJSONPath().Append("mobile1234"),
					ValMasked:   "12344321",
					KVFieldRel:  nil,
				},
				{
					Key:         "phonenumber",
					Val:         1234567890,
					ValMasked:   1234567890,
					ValJSONPath: types.NewJSONPath().Append("phonenumber"),
					KVFieldRel:  nil,
				},
				{
					Key:         "status",
					Val:         nil,
					ValMasked:   nil,
					ValJSONPath: types.NewJSONPath().Append("status"),
					KVFieldRel:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "match regex val",
			args: args{
				rule: types.KVRules{
					DetectRules: []types.KVDetectConfig{
						{
							ValRegex: []string{`^LTA[a-zA-Z0-9]+$`, `^1[0-9]+$`},
						},
					},
				},
				input: map[string]interface{}{
					"matchStr":    "LTAabcdEFGH1234",
					"matchInt":    112345,
					"mismatchStr": "ATAabcdEFGH1234",
					"mismatchInt": 212345,
				},
			},
			wantDetect: []types.KVPair{
				{
					Key:         "matchInt",
					Val:         112345,
					ValMasked:   112345,
					ValJSONPath: types.NewJSONPath().Append("matchInt"),
					KVFieldRel:  nil,
				},
				{
					Key:         "matchStr",
					Val:         "LTAabcdEFGH1234",
					ValMasked:   "LTAabcdEFGH1234",
					ValJSONPath: types.NewJSONPath().Append("matchStr"),
					KVFieldRel:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "match kv fields & matchMode=and",
			args: args{
				rule: types.KVRules{
					DetectRules: []types.KVDetectConfig{
						{
							KeyEqs: []string{
								"real_key",
							},
							ValEqs: []string{
								"real_val",
							},
							MatchMode: types.KVMatchAnd,
							KVFieldOpt: &types.KVField{
								Key: "find_key",
								Val: "find_val",
							},
						},
					},
				},
				input: map[string]interface{}{
					"find_key": "mismatch_real_key",
					"find_val": "real_val",
					"kv": map[string]string{
						"find_key": "real_key",
						"find_val": "real_val",
					},
					"kv2": map[string]interface{}{
						"find_key": "real_key",
						"find_val": []string{
							"real_val",
							"mismatch_real_val",
						},
					},
					"kv3": map[string]interface{}{
						"find_key": "mismatch_real_key",
						"find_val": []string{
							"real_val",
							"mismatch_real_val",
						},
					},
				},
			},
			wantDetect: []types.KVPair{
				{
					Key:         "real_key",
					Val:         "real_val",
					ValJSONPath: types.NewJSONPath().Append("kv").Append("find_val"),
					ValMasked:   "real_val",
					KVFieldRel: &types.KVField{
						Key: "find_key",
						Val: "find_val",
					},
				},
				{
					Key:         "real_key",
					Val:         "real_val",
					ValJSONPath: types.NewJSONPath().Append("kv2").Append("find_val").Append(0),
					ValMasked:   "real_val",
					KVFieldRel: &types.KVField{
						Key: "find_key",
						Val: "find_val",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJsonBytes, err := json.Marshal(tt.args.input)
			if err != nil {
				t.Errorf("json.Marshal() inputJsonBytes error = %v", err)
				return
			}
			wantDetectBytes, err := json.Marshal(tt.wantDetect)
			if err != nil {
				t.Errorf("json.Marshal() wantDetectBytes error = %v", err)
				return
			}
			m := NewKVProcesser(tt.args.rule)
			got, err := m.Detect(inputJsonBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Detect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotPairBytes, err := json.Marshal(got)
			if err != nil {
				t.Errorf("json.Marshal() gotPairBytes error = %v", err)
				return
			}
			if string(gotPairBytes) != string(wantDetectBytes) {
				t.Errorf("Detect() = %v, want %v", string(gotPairBytes), string(wantDetectBytes))
			}
		})
	}
}

func BenchmarkJSON_Detect(b *testing.B) {
	rule := types.KVRules{
		DetectRules: []types.KVDetectConfig{
			{
				KeyEqs: []string{
					"keyname",
				},
				ValEqs: []string{
					"value",
				},
				MatchMode: types.KVMatchAnd,
				KVFieldOpt: &types.KVField{
					Key: "key",
					Val: "val",
				},
			},
			{
				ValEqs: []string{"aaaa"},
			},
			{
				ValContains: []string{"919191"},
			},
			{
				ValEqs: []string{"null"},
			},
			{
				ValEqs: []string{"Value"},
			},
		},
	}
	input := `
	{
		"arr": [
			"aaaa",
			11111,
			true,
			null
		],
		"kv": {
			"key": "keyname",
			"kv": {
				"key": "keyname",
				"val": "value"
			},
			"val": "value"
		},
		"mobile": "+85291919191",
		"nullValue": null,
		"obj": {
			"mobile": "+85291919191",
			"nullValue": null,
			"password": "1234567890",
			"phone": 91919191,
			"username": "test"
		},
		"password": "1234567890",
		"phone": 91919191,
		"username": "test"
	}
	`
	inputByte := []byte(input)
	m := NewKVProcesser(rule)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Detect(inputByte)
	}
}

func BenchmarkVisitJSON(b *testing.B) {
	input := `
	{
		"arr": [
			"aaaa",
			11111,
			true,
			null
		],
		"kv": {
			"key": "keyname",
			"kv": {
				"key": "keyname",
				"val": "value"
			},
			"val": "value"
		},
		"mobile": "+85291919191",
		"nullValue": null,
		"obj": {
			"mobile": "+85291919191",
			"nullValue": null,
			"password": "1234567890",
			"phone": 91919191,
			"username": "test"
		},
		"password": "1234567890",
		"phone": 91919191,
		"username": "test"
	}
	`
	result, err := fastjson.Parse(input)
	if err != nil {
		b.Fatal(err)
	}
	m := NewKVProcesser(types.KVRules{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.visit(types.NewJSONPath(), "", result, nil)
	}
}

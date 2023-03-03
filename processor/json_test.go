package processer

import (
	"encoding/json"
	"testing"

	"github.com/senayuki/mosaic/types"
)

func TestMaskJSON(t *testing.T) {
	type args struct {
		rule  types.DLPRule
		input interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "mask json object",
			args: args{
				rule: types.DLPRule{},
				input: map[string]interface{}{
					"username": "test",
					"password": "1234567890",
					"mobile":   "+85291919191",
					"phone":    91919191,
					"arr": []interface{}{
						"aaaa",
						11111,
						true,
					},
					"nullValue": nil,
					"obj": map[string]interface{}{
						"username":  "test",
						"password":  "1234567890",
						"mobile":    "+85291919191",
						"phone":     91919191,
						"nullValue": nil,
					},
				},
			},
			want: map[string]interface{}{
				"username": "test",
				"password": "1234567890",
				"mobile":   "+85291919191",
				"phone":    91919191,
				"arr": []interface{}{
					"aaaa",
					11111,
					true,
				},
				"nullValue": nil,
				"obj": map[string]interface{}{
					"username":  "test",
					"password":  "1234567890",
					"mobile":    "+85291919191",
					"phone":     91919191,
					"nullValue": nil,
				},
			},
			wantErr: false,
		},
		{
			name: "mask json k-v fields object",
			args: args{
				rule: types.DLPRule{},
				input: map[string]interface{}{
					"customKey":  "password",
					"customVal":  "123456",
					"customVal2": 123456,
					"subObjHit": map[string]interface{}{
						"customKey":  "password",
						"customVal":  "123456",
						"customVal2": 123456,
					},
					"subObjMissKey": map[string]interface{}{
						"customKeyMiss": "password",
						"customVal":     "123456",
						"customVal2":    123456,
					},
					"subObjMissMatch": map[string]interface{}{
						"customKeyMiss": "mismatch",
						"customVal":     "123456",
						"customVal2":    123456,
					},
				},
			},
			want: map[string]interface{}{
				"customKey":  "password",
				"customVal":  "123456",
				"customVal2": 123456,
				"subObjHit": map[string]interface{}{
					"customKey":  "password",
					"customVal":  "123456",
					"customVal2": 123456,
				},
				"subObjMissKey": map[string]interface{}{
					"customKeyMiss": "password",
					"customVal":     "123456",
					"customVal2":    123456,
				},
				"subObjMissMatch": map[string]interface{}{
					"customKeyMiss": "mismatch",
					"customVal":     "123456",
					"customVal2":    123456,
				},
			},
			wantErr: false,
		},
		{
			name: "mask json array",
			args: args{
				rule: types.DLPRule{},
				input: []interface{}{
					"test-user",
					"123456",
					"+85291919191",
					91919191,
					nil,
				},
			},
			want: []interface{}{
				"test-user",
				"123456",
				"+85291919191",
				91919191,
				nil,
			},
			wantErr: false,
		},
		{
			name: "mask json null",
			args: args{
				rule:  types.DLPRule{},
				input: nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "mask json string",
			args: args{
				rule:  types.DLPRule{},
				input: "string",
			},
			want:    "string",
			wantErr: false,
		},
		{
			name: "mask json number",
			args: args{
				rule:  types.DLPRule{},
				input: 12121212,
			},
			want:    12121212,
			wantErr: false,
		},
		{
			name: "mask json bool",
			args: args{
				rule:  types.DLPRule{},
				input: true,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJsonBytes, err := json.Marshal(tt.args.input)
			wantJsonBytes, err := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}
			got, err := MaskJSON(tt.args.rule, string(inputJsonBytes))
			if (err != nil) != tt.wantErr {
				t.Errorf("MaskJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != string(wantJsonBytes) {
				t.Errorf("MaskJSON() = %v, want %v", got, string(wantJsonBytes))
			}
		})
	}
}

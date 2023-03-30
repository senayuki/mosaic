package mask

import (
	"testing"

	"github.com/senayuki/mosaic/types"
)

func TestMarkCoverProcesser_Mask(t *testing.T) {
	type args struct {
		in string
		kv *types.KVPair
	}
	tests := []struct {
		name    string
		fields  types.KVMaskConfig
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "empty",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Length:  3,
					Reverse: false,
				},
			},
			args: args{
				in: "",
			},
			wantOut: "",
			wantErr: false,
		},
		{
			name: "normal mask",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Length:  3,
					Reverse: false,
				},
			},
			args: args{
				in: "I can eat glass, it does not hurt me",
			},
			wantOut: "***",
			wantErr: false,
		},
		{
			name: "normal mask with offset",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Char:    "*",
					Offset:  3,
					Padding: 3,
					Length:  3,
					Reverse: false,
				},
			},
			args: args{
				in: "I can eat glass, it does not hurt me",
			},
			wantOut: "I c*** me",
			wantErr: false,
		},
		{
			name: "CJK mask",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Char:    "*",
					Offset:  3,
					Padding: 3,
					Length:  3,
					Reverse: false,
				},
			},
			args: args{
				in: "我能吞下玻璃而不伤身体",
			},
			wantOut: "我能吞***伤身体",
			wantErr: false,
		},
		{
			name: "offset more than content",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Char:    "*",
					Offset:  20,
					Reverse: false,
				},
			},
			args: args{
				in: "我能吞下玻璃而不伤身体",
			},
			wantOut: "我能吞下玻璃而不伤身*",
			wantErr: false,
		},
		{
			name: "padding more than content",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Char:    "*",
					Padding: 20,
					Reverse: false,
				},
			},
			args: args{
				in: "我能吞下玻璃而不伤身体",
			},
			wantOut: "*能吞下玻璃而不伤身体",
			wantErr: false,
		},
		{
			name: "offset & padding more than content",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Char:    "*",
					Offset:  20,
					Padding: 20,
					Reverse: false,
				},
			},
			args: args{
				in: "我能吞下玻璃而不伤身体",
			},
			wantOut: "我能吞下玻*而不伤身体",
			wantErr: false,
		},
		{
			name: "offset & padding more than content",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Char:    "*",
					Offset:  20,
					Padding: 20,
					Reverse: false,
				},
			},
			args: args{
				in: "我能吞下玻璃而不伤身体啊",
			},
			wantOut: "我能吞下玻**不伤身体啊",
			wantErr: false,
		},
		{
			name: "offset & padding single char",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Char:    "*",
					Offset:  20,
					Padding: 20,
					Reverse: false,
				},
			},
			args: args{
				in: "I",
			},
			wantOut: "*",
			wantErr: false,
		},
		{
			name: "length more than content",
			fields: types.KVMaskConfig{
				CoverParam: types.MaskRuleCoverParam{
					Char:    "*",
					Offset:  3,
					Padding: 3,
					Length:  20,
					Reverse: false,
				},
			},
			args: args{
				in: "我能吞下玻璃而不伤身体",
			},
			wantOut: "我能吞*****伤身体",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MarkCoverProcesser{}
			m.Init(&tt.fields)
			gotOut, err := m.Mask(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkCoverProcesser.Mask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut != tt.wantOut {
				t.Errorf("MarkCoverProcesser.Mask() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

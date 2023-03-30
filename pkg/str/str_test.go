package str

import (
	"reflect"
	"testing"
)

func TestRunes2Bytes(t *testing.T) {
	type args struct {
		in []rune
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{
				in: []rune("Hello 世界"),
			},
			want: []byte("Hello 世界"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Runes2Bytes(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Runes2Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes2Runes(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name string
		args args
		want []rune
	}{
		{
			name: "normal",
			args: args{
				in: []byte("Hello 世界"),
			},
			want: []rune("Hello 世界"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bytes2Runes(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes2Runes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkRunes2Bytes(b *testing.B) {
	runes := []rune("Hello 世界")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Runes2Bytes(runes)
	}
}

func BenchmarkBytes2Runes(b *testing.B) {
	bytes := []byte("Hello 世界")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Bytes2Runes(bytes)
	}
}

package json

import (
	"math"
	"reflect"
	"testing"
)

func Test_appendFloat64(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  []byte
	}{
		{"-Inf", math.Inf(-1), []byte(`"-Inf"`)},
		{"+Inf", math.Inf(1), []byte(`"+Inf"`)},
		{"NaN", math.NaN(), []byte(`"NaN"`)},
		{"0", 0, []byte(`0`)},
		{"-1.1", -1.1, []byte(`-1.1`)},
		{"1e20", 1e20, []byte(`100000000000000000000`)},
		{"1e21", 1e21, []byte(`1000000000000000000000`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendFloat32([]byte{}, float32(tt.input)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendFloat32() = %s, want %s", got, tt.want)
			}
			if got := AppendFloat64([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendFloat32() = %s, want %s", got, tt.want)
			}
		})
	}
}

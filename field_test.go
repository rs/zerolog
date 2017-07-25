package zerolog

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
		{"-Inf", math.Inf(-1), []byte(`"foo":"-Inf"`)},
		{"+Inf", math.Inf(1), []byte(`"foo":"+Inf"`)},
		{"NaN", math.NaN(), []byte(`"foo":"NaN"`)},
		{"0", 0, []byte(`"foo":0`)},
		{"-1.1", -1.1, []byte(`"foo":-1.1`)},
		{"1e20", 1e20, []byte(`"foo":100000000000000000000`)},
		{"1e21", 1e21, []byte(`"foo":1000000000000000000000`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendFloat32([]byte{}, "foo", float32(tt.input)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendFloat32() = %s, want %s", got, tt.want)
			}
			if got := appendFloat64([]byte{}, "foo", tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendFloat32() = %s, want %s", got, tt.want)
			}
		})
	}
}

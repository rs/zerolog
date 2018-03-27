package json

import (
	"math"
	"reflect"
	"testing"
)

func TestAppendType(t *testing.T) {
	w := map[string]func(interface{}) []byte{
		"AppendInt":     func(v interface{}) []byte { return AppendInt([]byte{}, v.(int)) },
		"AppendInt8":    func(v interface{}) []byte { return AppendInt8([]byte{}, v.(int8)) },
		"AppendInt16":   func(v interface{}) []byte { return AppendInt16([]byte{}, v.(int16)) },
		"AppendInt32":   func(v interface{}) []byte { return AppendInt32([]byte{}, v.(int32)) },
		"AppendInt64":   func(v interface{}) []byte { return AppendInt64([]byte{}, v.(int64)) },
		"AppendUint":    func(v interface{}) []byte { return AppendUint([]byte{}, v.(uint)) },
		"AppendUint8":   func(v interface{}) []byte { return AppendUint8([]byte{}, v.(uint8)) },
		"AppendUint16":  func(v interface{}) []byte { return AppendUint16([]byte{}, v.(uint16)) },
		"AppendUint32":  func(v interface{}) []byte { return AppendUint32([]byte{}, v.(uint32)) },
		"AppendUint64":  func(v interface{}) []byte { return AppendUint64([]byte{}, v.(uint64)) },
		"AppendFloat32": func(v interface{}) []byte { return AppendFloat32([]byte{}, v.(float32)) },
		"AppendFloat64": func(v interface{}) []byte { return AppendFloat64([]byte{}, v.(float64)) },
	}
	tests := []struct {
		name  string
		fn    string
		input interface{}
		want  []byte
	}{
		{"AppendInt8(math.MaxInt8)", "AppendInt8", int8(math.MaxInt8), []byte("127")},
		{"AppendInt16(math.MaxInt16)", "AppendInt16", int16(math.MaxInt16), []byte("32767")},
		{"AppendInt32(math.MaxInt32)", "AppendInt32", int32(math.MaxInt32), []byte("2147483647")},
		{"AppendInt64(math.MaxInt64)", "AppendInt64", int64(math.MaxInt64), []byte("9223372036854775807")},

		{"AppendUint8(math.MaxUint8)", "AppendUint8", uint8(math.MaxUint8), []byte("255")},
		{"AppendUint16(math.MaxUint16)", "AppendUint16", uint16(math.MaxUint16), []byte("65535")},
		{"AppendUint32(math.MaxUint32)", "AppendUint32", uint32(math.MaxUint32), []byte("4294967295")},
		{"AppendUint64(math.MaxUint64)", "AppendUint64", uint64(math.MaxUint64), []byte("18446744073709551615")},

		{"AppendFloat32(-Inf)", "AppendFloat32", float32(math.Inf(-1)), []byte(`"-Inf"`)},
		{"AppendFloat32(+Inf)", "AppendFloat32", float32(math.Inf(1)), []byte(`"+Inf"`)},
		{"AppendFloat32(NaN)", "AppendFloat32", float32(math.NaN()), []byte(`"NaN"`)},
		{"AppendFloat32(0)", "AppendFloat32", float32(0), []byte(`0`)},
		{"AppendFloat32(-1.1)", "AppendFloat32", float32(-1.1), []byte(`-1.1`)},
		{"AppendFloat32(1e20)", "AppendFloat32", float32(1e20), []byte(`100000000000000000000`)},
		{"AppendFloat32(1e21)", "AppendFloat32", float32(1e21), []byte(`1000000000000000000000`)},

		{"AppendFloat64(-Inf)", "AppendFloat64", float64(math.Inf(-1)), []byte(`"-Inf"`)},
		{"AppendFloat64(+Inf)", "AppendFloat64", float64(math.Inf(1)), []byte(`"+Inf"`)},
		{"AppendFloat64(NaN)", "AppendFloat64", float64(math.NaN()), []byte(`"NaN"`)},
		{"AppendFloat64(0)", "AppendFloat64", float64(0), []byte(`0`)},
		{"AppendFloat64(-1.1)", "AppendFloat64", float64(-1.1), []byte(`-1.1`)},
		{"AppendFloat64(1e20)", "AppendFloat64", float64(1e20), []byte(`100000000000000000000`)},
		{"AppendFloat64(1e21)", "AppendFloat64", float64(1e21), []byte(`1000000000000000000000`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := w[tt.fn](tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

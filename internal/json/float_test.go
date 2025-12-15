package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/rs/zerolog/internal"
)

var float64Tests = []struct {
	Name string
	Val  float64
	Want string
}{
	{
		Name: "Positive integer",
		Val:  1234.0,
		Want: "1234",
	},
	{
		Name: "Negative integer",
		Val:  -5678.0,
		Want: "-5678",
	},
	{
		Name: "Positive decimal",
		Val:  12.3456,
		Want: "12.3456",
	},
	{
		Name: "Negative decimal",
		Val:  -78.9012,
		Want: "-78.9012",
	},
	{
		Name: "Large positive number",
		Val:  123456789.0,
		Want: "123456789",
	},
	{
		Name: "Large negative number",
		Val:  -987654321.0,
		Want: "-987654321",
	},
	{
		Name: "Zero",
		Val:  0.0,
		Want: "0",
	},
	{
		Name: "Smallest positive value",
		Val:  math.SmallestNonzeroFloat64,
		Want: "5e-324",
	},
	{
		Name: "Largest positive value",
		Val:  math.MaxFloat64,
		Want: "1.7976931348623157e+308",
	},
	{
		Name: "Smallest negative value",
		Val:  -math.SmallestNonzeroFloat64,
		Want: "-5e-324",
	},
	{
		Name: "Largest negative value",
		Val:  -math.MaxFloat64,
		Want: "-1.7976931348623157e+308",
	},
	{
		Name: "NaN",
		Val:  math.NaN(),
		Want: `"NaN"`,
	},
	{
		Name: "+Inf",
		Val:  math.Inf(1),
		Want: `"+Inf"`,
	},
	{
		Name: "-Inf",
		Val:  math.Inf(-1),
		Want: `"-Inf"`,
	},
	{
		Name: "Clean up e-09 to e-9 case 1",
		Val:  1e-9,
		Want: "1e-9",
	},
	{
		Name: "Clean up e-09 to e-9 case 2",
		Val:  -2.236734e-9,
		Want: "-2.236734e-9",
	},
}

func TestEncoder_AppendFloat64(t *testing.T) {
	for _, tc := range float64Tests {
		t.Run(tc.Name, func(t *testing.T) {
			var b []byte
			b = (Encoder{}).AppendFloat64(b, tc.Val, -1)
			if s := string(b); tc.Want != s {
				t.Errorf("%q", s)
			}
		})
	}
}

func FuzzEncoder_AppendFloat64(f *testing.F) {
	for _, tc := range float64Tests {
		f.Add(tc.Val)
	}
	f.Fuzz(func(t *testing.T, val float64) {
		actual := (Encoder{}).AppendFloat64(nil, val, -1)
		if len(actual) == 0 {
			t.Fatal("empty buffer")
		}

		if actual[0] == '"' {
			switch string(actual) {
			case `"NaN"`:
				if !math.IsNaN(val) {
					t.Fatalf("expected %v got NaN", val)
				}
			case `"+Inf"`:
				if !math.IsInf(val, 1) {
					t.Fatalf("expected %v got +Inf", val)
				}
			case `"-Inf"`:
				if !math.IsInf(val, -1) {
					t.Fatalf("expected %v got -Inf", val)
				}
			default:
				t.Fatalf("unexpected string: %s", actual)
			}
			return
		}

		if expected, err := json.Marshal(val); err != nil {
			t.Error(err)
		} else if string(actual) != string(expected) {
			t.Errorf("expected %s, got %s", expected, actual)
		}

		var parsed float64
		if err := json.Unmarshal(actual, &parsed); err != nil {
			t.Fatal(err)
		}

		if parsed != val {
			t.Fatalf("expected %v, got %v", val, parsed)
		}
	})
}

var float32Tests = []struct {
	Name string
	Val  float32
	Want string
}{
	{
		Name: "Positive integer",
		Val:  1234.0,
		Want: "1234",
	},
	{
		Name: "Negative integer",
		Val:  -5678.0,
		Want: "-5678",
	},
	{
		Name: "Positive decimal",
		Val:  12.3456,
		Want: "12.3456",
	},
	{
		Name: "Negative decimal",
		Val:  -78.9012,
		Want: "-78.9012",
	},
	{
		Name: "Large positive number",
		Val:  123456789.0,
		Want: "123456790",
	},
	{
		Name: "Large negative number",
		Val:  -987654321.0,
		Want: "-987654340",
	},
	{
		Name: "Zero",
		Val:  0.0,
		Want: "0",
	},
	{
		Name: "Smallest positive value",
		Val:  math.SmallestNonzeroFloat32,
		Want: "1e-45",
	},
	{
		Name: "Largest positive value",
		Val:  math.MaxFloat32,
		Want: "3.4028235e+38",
	},
	{
		Name: "Smallest negative value",
		Val:  -math.SmallestNonzeroFloat32,
		Want: "-1e-45",
	},
	{
		Name: "Largest negative value",
		Val:  -math.MaxFloat32,
		Want: "-3.4028235e+38",
	},
	{
		Name: "NaN",
		Val:  float32(math.NaN()),
		Want: `"NaN"`,
	},
	{
		Name: "+Inf",
		Val:  float32(math.Inf(1)),
		Want: `"+Inf"`,
	},
	{
		Name: "-Inf",
		Val:  float32(math.Inf(-1)),
		Want: `"-Inf"`,
	},
	{
		Name: "Clean up e-09 to e-9 case 1",
		Val:  1e-9,
		Want: "1e-9",
	},
	{
		Name: "Clean up e-09 to e-9 case 2",
		Val:  -2.236734e-9,
		Want: "-2.236734e-9",
	},
}

func TestEncoder_AppendFloat32(t *testing.T) {
	for _, tc := range float32Tests {
		t.Run(tc.Name, func(t *testing.T) {
			var b []byte
			b = (Encoder{}).AppendFloat32(b, tc.Val, -1)
			if s := string(b); tc.Want != s {
				t.Errorf("%q", s)
			}
		})
	}
}

func FuzzEncoder_AppendFloat32(f *testing.F) {
	for _, tc := range float32Tests {
		f.Add(tc.Val)
	}
	f.Fuzz(func(t *testing.T, val float32) {
		actual := (Encoder{}).AppendFloat32(nil, val, -1)
		if len(actual) == 0 {
			t.Fatal("empty buffer")
		}

		if actual[0] == '"' {
			val := float64(val)
			switch string(actual) {
			case `"NaN"`:
				if !math.IsNaN(val) {
					t.Fatalf("expected %v got NaN", val)
				}
			case `"+Inf"`:
				if !math.IsInf(val, 1) {
					t.Fatalf("expected %v got +Inf", val)
				}
			case `"-Inf"`:
				if !math.IsInf(val, -1) {
					t.Fatalf("expected %v got -Inf", val)
				}
			default:
				t.Fatalf("unexpected string: %s", actual)
			}
			return
		}

		if expected, err := json.Marshal(val); err != nil {
			t.Error(err)
		} else if string(actual) != string(expected) {
			t.Errorf("expected %s, got %s", expected, actual)
		}

		var parsed float32
		if err := json.Unmarshal(actual, &parsed); err != nil {
			t.Fatal(err)
		}

		if parsed != val {
			t.Fatalf("expected %v, got %v", val, parsed)
		}
	})
}

func generateFloat32s(n int) []float32 {
	floats := make([]float32, n)
	for i := 0; i < n; i++ {
		floats[i] = rand.Float32()
	}
	return floats
}

func generateFloat64s(n int) []float64 {
	floats := make([]float64, n)
	for i := 0; i < n; i++ {
		floats[i] = rand.Float64()
	}
	return floats
}

func TestAppendFloats32(t *testing.T) {
	doOne := func(vals []float32) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			if math.IsNaN(float64(val)) {
				want = append(want, []byte(`"NaN"`)...)
			} else if math.IsInf(float64(val), 1) {
				want = append(want, []byte(`"+Inf"`)...)
			} else if math.IsInf(float64(val), -1) {
				want = append(want, []byte(`"-Inf"`)...)
			} else {
				want = append(want, []byte(fmt.Sprintf("%v", float32(val)))...)
			}
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendFloats32([]byte{}, vals, -1)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendFloats32(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]float32, 0)
	for _, tc := range internal.Float32TestCases {
		if tc.Val > 0 && tc.Val < 1e-4 {
			continue // we want to ignore very small numbers for this test
		}
		array = append(array, float32(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}

func TestAppendFloats64(t *testing.T) {
	doOne := func(vals []float64) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			if math.IsNaN(val) {
				want = append(want, []byte(`"NaN"`)...)
			} else if math.IsInf(val, 1) {
				want = append(want, []byte(`"+Inf"`)...)
			} else if math.IsInf(val, -1) {
				want = append(want, []byte(`"-Inf"`)...)
			} else {
				want = append(want, []byte(fmt.Sprintf("%v", val))...)
			}
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendFloats64([]byte{}, vals, -1)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendFloats64(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]float64, 0)
	for _, tc := range internal.Float64TestCases {
		if tc.Val > 0 && tc.Val < 1e-4 {
			continue // we want to ignore very small numbers for this test
		}
		array = append(array, tc.Val)
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}

// this is really just for the memory allocation characteristics
func BenchmarkEncoder_AppendFloat32(b *testing.B) {
	floats := append(generateFloat32s(5000), float32(math.NaN()), float32(math.Inf(1)), float32(math.Inf(-1)))
	dst := make([]byte, 0, 128)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, f := range floats {
			dst = (Encoder{}).AppendFloat32(dst[:0], f, -1)
		}
	}
}

// this is really just for the memory allocation characteristics
func BenchmarkEncoder_AppendFloat64(b *testing.B) {
	floats := append(generateFloat64s(5000), math.NaN(), math.Inf(1), math.Inf(-1))
	dst := make([]byte, 0, 128)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, f := range floats {
			dst = (Encoder{}).AppendFloat64(dst[:0], f, -1)
		}
	}
}

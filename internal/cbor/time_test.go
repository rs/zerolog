package cbor

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/internal"
)

func TestEncoder_AppendDuration(t *testing.T) {
	type args struct {
		dst    []byte
		d      time.Duration
		unit   time.Duration
		format string
		useInt bool
		unused int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "useInt",
			args: args{
				d:      1234567890,
				unit:   time.Second,
				useInt: true,
			},
			want: []byte{1},
		},
		{
			name: "formatFloat",
			args: args{
				d:      1234567890,
				unit:   time.Second,
				format: durationFormatFloat,
			},
			want: []byte{251, 63, 243, 192, 202, 66, 131, 222, 27},
		},
		{
			name: "formatInt",
			args: args{
				d:      1234567890,
				unit:   time.Second,
				format: durationFormatInt,
			},
			want: []byte{1},
		},
		{
			name: "formatString",
			args: args{
				d:      1234567890,
				unit:   time.Second,
				format: durationFormatString,
			},
			want: []byte{107, 49, 46, 50, 51, 52, 53, 54, 55, 56, 57, 115},
		},
		{
			name: "formatBlank",
			args: args{
				d:    1234567890,
				unit: time.Second,
			},
			want: []byte{251, 63, 243, 192, 202, 66, 131, 222, 27},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Encoder{}
			if got := e.AppendDuration(tt.args.dst, tt.args.d, tt.args.unit, tt.args.format, tt.args.useInt, tt.args.unused); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoder_AppendDurations(t *testing.T) {
	type args struct {
		dst    []byte
		vals   []time.Duration
		unit   time.Duration
		format string
		useInt bool
		unused int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "useInt",
			args: args{
				vals:   []time.Duration{1234567890},
				unit:   time.Second,
				useInt: true,
			},
			want: []byte{129, 1},
		},
		{
			name: "formatFloat",
			args: args{
				vals:   []time.Duration{1234567890},
				unit:   time.Second,
				format: durationFormatFloat,
			},
			want: []byte{129, 251, 63, 243, 192, 202, 66, 131, 222, 27},
		},
		{
			name: "formatInt",
			args: args{
				vals:   []time.Duration{1234567890},
				unit:   time.Second,
				format: durationFormatInt,
			},
			want: []byte{129, 1},
		},
		{
			name: "formatString",
			args: args{
				vals:   []time.Duration{1234567890},
				unit:   time.Second,
				format: durationFormatString,
			},
			want: []byte{129, 107, 49, 46, 50, 51, 52, 53, 54, 55, 56, 57, 115},
		},
		{
			name: "formatBlank",
			args: args{
				vals: []time.Duration{1234567890},
				unit: time.Second,
			},
			want: []byte{129, 251, 63, 243, 192, 202, 66, 131, 222, 27},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Encoder{}
			if got := e.AppendDurations(tt.args.dst, tt.args.vals, tt.args.unit, tt.args.format, tt.args.useInt, tt.args.unused); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendDurations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppendTimeNow(t *testing.T) {
	tm := time.Now()
	s := enc.AppendTime([]byte{}, tm, "unused")
	got := string(s)

	tm1 := float64(tm.Unix()) + float64(tm.Nanosecond())*1e-9
	tm2 := math.Float64bits(tm1)
	var tm3 [8]byte
	for i := uint(0); i < 8; i++ {
		tm3[i] = byte(tm2 >> ((8 - i - 1) * 8))
	}
	want := append([]byte{0xc1, 0xfb}, tm3[:]...)
	if got != string(want) {
		t.Errorf("Appendtime(%s)=0x%s, want: 0x%s",
			"time.Now()", hex.EncodeToString(s),
			hex.EncodeToString(want))
	}
}

func TestAppendTimePastPresentInteger(t *testing.T) {
	for _, tt := range internal.TimeIntegerTestcases {
		tin, err := time.Parse(time.RFC3339, tt.Txt)
		if err != nil {
			fmt.Println("Cannot parse input", tt.Txt, ".. Skipping!", err)
			continue
		}
		b := enc.AppendTime([]byte{}, tin, "unused")
		if got, want := string(b), tt.Binary; got != want {
			t.Errorf("appendString(%s) = 0x%s, want 0x%s", tt.Txt,
				hex.EncodeToString(b),
				hex.EncodeToString([]byte(want)))
		}
	}
}

func TestAppendTimePastPresentFloat(t *testing.T) {
	const timeFloatFmt = "2006-01-02T15:04:05.999999-07:00"
	for _, tt := range internal.TimeFloatTestcases {
		tin, err := time.Parse(timeFloatFmt, tt.RfcStr)
		if err != nil {
			fmt.Println("Cannot parse input", tt.RfcStr, ".. Skipping!")
			continue
		}
		b := enc.AppendTime([]byte{}, tin, "unused")
		if got, want := string(b), tt.Out; got != want {
			t.Errorf("appendString(%s) = 0x%s, want 0x%s", tt.RfcStr,
				hex.EncodeToString(b),
				hex.EncodeToString([]byte(want)))
		}
	}
}
func TestAppendTimes(t *testing.T) {
	const timeFloatFmt = "2006-01-02T15:04:05.999999-07:00"
	array := make([]time.Time, len(internal.TimeFloatTestcases))
	want := make([]byte, 0)
	want = append(want, 0x82) // start small array
	for i, tt := range internal.TimeFloatTestcases {
		array[i], _ = time.Parse(timeFloatFmt, tt.RfcStr)
		want = append(want, []byte(tt.Out)...)
	}

	got := enc.AppendTimes([]byte{}, array, "unused")
	if !bytes.Equal(got, want) {
		t.Errorf("AppendTimes(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]time.Time, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendTimes([]byte{}, array, "unused")
	if !bytes.Equal(got, want) {
		t.Errorf("AppendTimes(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now large array case
	testtime, _ := time.Parse(timeFloatFmt, internal.TimeFloatTestcases[0].RfcStr)
	outbytes := internal.TimeFloatTestcases[0].Out
	array = make([]time.Time, 24)
	want = make([]byte, 0)
	want = append(want, 0x98) // start a large array
	want = append(want, 0x18) // of length 24
	for i := 0; i < len(array); i++ {
		array[i] = testtime
		want = append(want, []byte(outbytes)...)
	}
	got = enc.AppendTimes([]byte{}, array, "unused")
	if !bytes.Equal(got, want) {
		t.Errorf("AppendTimes(%v)\ngot:  0x%s\nwant: 0x%s",
			array,
			hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendDurationFloat(t *testing.T) {
	for _, tt := range internal.DurTestcases {
		dur := tt.Duration
		want := []byte{}
		want = append(want, []byte(tt.FloatOut)...)
		got := enc.AppendDuration([]byte{}, dur, time.Microsecond, "", false, -1)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendDuration(%v)=\ngot:  0x%s\nwant: 0x%s",
				dur,
				hex.EncodeToString(got),
				hex.EncodeToString(want))
		}
	}
}
func TestAppendDurationInteger(t *testing.T) {
	for _, tt := range internal.DurTestcases {
		dur := tt.Duration
		want := []byte{}
		want = append(want, []byte(tt.IntegerOut)...)
		got := enc.AppendDuration([]byte{}, dur, time.Microsecond, "", true, -1)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendDuration(%v)=\ngot:  0x%s\nwant: 0x%s",
				dur,
				hex.EncodeToString(got),
				hex.EncodeToString(want))
		}
	}
}
func TestAppendDurations(t *testing.T) {
	array := make([]time.Duration, len(internal.DurTestcases))
	want := make([]byte, 0)
	want = append(want, 0x83) // start 3 element array
	for i, tt := range internal.DurTestcases {
		array[i] = tt.Duration
		want = append(want, []byte(tt.FloatOut)...)
	}

	got := enc.AppendDurations([]byte{}, array, time.Microsecond, "", false, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendDurations(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]time.Duration, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendDurations([]byte{}, array, time.Microsecond, "", false, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendDurations(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now large array case
	testtime := internal.DurTestcases[0].Duration
	outbytes := internal.DurTestcases[0].FloatOut
	array = make([]time.Duration, 24)
	want = make([]byte, 0)
	want = append(want, 0x98) // start a large array
	want = append(want, 0x18) // of length 24
	for i := 0; i < len(array); i++ {
		array[i] = testtime
		want = append(want, []byte(outbytes)...)
	}
	got = enc.AppendDurations([]byte{}, array, time.Microsecond, "", false, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendDurations(%v)\ngot:  0x%s\nwant: 0x%s",
			array,
			hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func BenchmarkAppendTime(b *testing.B) {
	tests := map[string]string{
		"Integer": "Feb 3, 2013 at 7:54pm (PST)",
		"Float":   "2006-01-02T15:04:05.999999-08:00",
	}
	const timeFloatFmt = "2006-01-02T15:04:05.999999-07:00"

	for name, str := range tests {
		t, err := time.Parse(time.RFC3339, str)
		if err != nil {
			t, _ = time.Parse(timeFloatFmt, str)
		}
		b.Run(name, func(b *testing.B) {
			buf := make([]byte, 0, 100)
			for i := 0; i < b.N; i++ {
				_ = enc.AppendTime(buf, t, "unused")
			}
		})
	}
}

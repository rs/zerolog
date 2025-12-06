package cbor

import (
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

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

var timeIntegerTestcases = []struct {
	txt    string
	binary string
	rfcStr string
}{
	{"2013-02-03T19:54:00-08:00", "\xc1\x1a\x51\x0f\x30\xd8", "2013-02-04T03:54:00Z"},
	{"1950-02-03T19:54:00-08:00", "\xc1\x3a\x25\x71\x93\xa7", "1950-02-04T03:54:00Z"},
}

func TestAppendTimePastPresentInteger(t *testing.T) {
	for _, tt := range timeIntegerTestcases {
		tin, err := time.Parse(time.RFC3339, tt.txt)
		if err != nil {
			fmt.Println("Cannot parse input", tt.txt, ".. Skipping!", err)
			continue
		}
		b := enc.AppendTime([]byte{}, tin, "unused")
		if got, want := string(b), tt.binary; got != want {
			t.Errorf("appendString(%s) = 0x%s, want 0x%s", tt.txt,
				hex.EncodeToString(b),
				hex.EncodeToString([]byte(want)))
		}
	}
}

var timeFloatTestcases = []struct {
	rfcStr string
	out    string
}{
	{"2006-01-02T15:04:05.999999-08:00", "\xc1\xfb\x41\xd0\xee\x6c\x59\x7f\xff\xfc"},
	{"1956-01-02T15:04:05.999999-08:00", "\xc1\xfb\xc1\xba\x53\x81\x1a\x00\x00\x11"},
}

func TestAppendTimePastPresentFloat(t *testing.T) {
	const timeFloatFmt = "2006-01-02T15:04:05.999999-07:00"
	for _, tt := range timeFloatTestcases {
		tin, err := time.Parse(timeFloatFmt, tt.rfcStr)
		if err != nil {
			fmt.Println("Cannot parse input", tt.rfcStr, ".. Skipping!")
			continue
		}
		b := enc.AppendTime([]byte{}, tin, "unused")
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendString(%s) = 0x%s, want 0x%s", tt.rfcStr,
				hex.EncodeToString(b),
				hex.EncodeToString([]byte(want)))
		}
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

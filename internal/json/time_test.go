package json

import (
	"bytes"
	"fmt"
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
			want: []byte{49},
		},
		{
			name: "formatFloat",
			args: args{
				d:      1234567890,
				unit:   time.Second,
				format: durationFormatFloat,
			},
			want: []byte{49},
		},
		{
			name: "formatInt",
			args: args{
				d:      1234567890,
				unit:   time.Second,
				format: durationFormatInt,
			},
			want: []byte{49},
		},
		{
			name: "formatString",
			args: args{
				d:      1234567890,
				unit:   time.Second,
				format: durationFormatString,
			},
			want: []byte{34, 49, 46, 50, 51, 52, 53, 54, 55, 56, 57, 115, 34},
		},
		{
			name: "formatBlank",
			args: args{
				d:    1234567890,
				unit: time.Second,
			},
			want: []byte{49},
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
			want: []byte{91, 49, 93},
		},
		{
			name: "formatFloat",
			args: args{
				vals:   []time.Duration{1234567890},
				unit:   time.Second,
				format: durationFormatFloat,
			},
			want: []byte{91, 49, 93},
		},
		{
			name: "formatInt",
			args: args{
				vals:   []time.Duration{1234567890},
				unit:   time.Second,
				format: durationFormatInt,
			},
			want: []byte{91, 49, 93},
		},
		{
			name: "formatString",
			args: args{
				vals:   []time.Duration{1234567890},
				unit:   time.Second,
				format: durationFormatString,
			},
			want: []byte{91, 34, 49, 46, 50, 51, 52, 53, 54, 55, 56, 57, 115, 34, 93},
		},
		{
			name: "formatBlank",
			args: args{
				vals: []time.Duration{1234567890},
				unit: time.Second,
			},
			want: []byte{91, 49, 93},
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
	got := enc.AppendTime([]byte{}, tm, time.RFC3339)
	want := tm.AppendFormat([]byte{'"'}, time.RFC3339)
	want = append(want, '"')
	if !bytes.Equal(got, want) {
		t.Errorf("AppendTime(%s)\ngot:  %s\nwant: %s",
			"time.Now()",
			string(got),
			string(want))
	}
}

func TestAppendTimePastPresentInteger(t *testing.T) {
	for _, tt := range internal.TimeIntegerTestcases {
		tin, err := time.Parse(time.RFC3339, tt.Txt)
		if err != nil {
			fmt.Println("Cannot parse input", tt.Txt, ".. Skipping!", err)
			continue
		}

		got := enc.AppendTime([]byte{}, tin, timeFormatUnix)
		want := []byte(fmt.Sprintf("%d", tt.UnixInt))
		if !bytes.Equal(got, want) {
			t.Errorf("appendString(%s)\ngot:  %s\nwant: %s",
				tt.Txt,
				string(got),
				string(want))
		}
		got = enc.AppendTime([]byte{}, tin, timeFormatUnixMs)
		want = []byte(fmt.Sprintf("%d", tt.UnixInt*1000))
		if !bytes.Equal(got, want) {
			t.Errorf("appendString(%s)\ngot:  %s\nwant: %s",
				tt.Txt,
				string(got),
				string(want))
		}
		got = enc.AppendTime([]byte{}, tin, timeFormatUnixMicro)
		want = []byte(fmt.Sprintf("%d", tt.UnixInt*1000000))
		if !bytes.Equal(got, want) {
			t.Errorf("appendString(%s)\ngot:  %s\nwant: %s",
				tt.Txt,
				string(got),
				string(want))
		}
		got = enc.AppendTime([]byte{}, tin, timeFormatUnixNano)
		want = []byte(fmt.Sprintf("%d", tt.UnixInt*1000000000))
		if !bytes.Equal(got, want) {
			t.Errorf("appendString(%s)\ngot:  %s\nwant: %s",
				tt.Txt,
				string(got),
				string(want))
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
		got := enc.AppendTime([]byte{}, tin, timeFormatUnix)
		want := []byte(fmt.Sprintf("%d", tt.UnixInt))
		if !bytes.Equal(got, want) {
			t.Errorf("appendString(%s)\ngot:  %s\nwant: %s",
				tt.RfcStr,
				string(got),
				string(want))
		}
	}
}
func TestAppendTimes(t *testing.T) {
	doOne := func(multiplier int, format string) {
		array := make([]time.Time, 0)
		want := append([]byte{}, '[')
		want = append(want, ']')
		got := enc.AppendTimes([]byte{}, array, format)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendTimes(%v)\ngot:  %s\nwant: %s",
				array,
				string(got),
				string(want))
		}

		array = make([]time.Time, len(internal.TimeIntegerTestcases))
		want = append([]byte{}, '[')
		for i, tt := range internal.TimeIntegerTestcases {
			if tin, err := time.Parse(time.RFC3339, tt.RfcStr); err != nil {
				fmt.Println("Cannot parse input", tt.RfcStr, ".. Skipping!")
				continue
			} else {
				array[i] = tin
			}
			if multiplier == 0 {
				want = append(want, '"')
				formatted := array[i].Format(format)
				want = append(want, []byte(fmt.Sprintf("%v", formatted))...)
				want = append(want, '"')
			} else {
				scaled := tt.UnixInt * multiplier
				want = append(want, []byte(fmt.Sprintf("%d", scaled))...)
			}
			if i < len(internal.TimeIntegerTestcases)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got = enc.AppendTimes([]byte{}, array, format)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendTimes(%v) %d %s\ngot:  %s\nwant: %s",
				array, multiplier, format,
				string(got),
				string(want))
		}
	}

	doOne(0, time.RFC3339)
	doOne(1, timeFormatUnix)
	doOne(1000, timeFormatUnixMs)
	doOne(1000000, timeFormatUnixMicro)
	doOne(1000000000, timeFormatUnixNano)
}

func TestAppendDurationFloat(t *testing.T) {
	for _, tt := range internal.DurTestcases {
		dur := tt.Duration
		want := []byte{}
		want = append(want, []byte(fmt.Sprintf("%v", dur.Microseconds()))...)
		got := enc.AppendDuration([]byte{}, dur, time.Microsecond, "", false, -1)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendDuration(%v)=\ngot:  %s\nwant: %s",
				dur,
				string(got),
				string(want))
		}

		want = []byte{}
		fraction := float64(dur) / float64(time.Millisecond)
		want = append(want, []byte(fmt.Sprintf("%v", fraction))...)
		got = enc.AppendDuration([]byte{}, dur, time.Millisecond, "", false, -1)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendDuration(%v)=\ngot:  %s\nwant: %s",
				dur,
				string(got),
				string(want))
		}
	}
}
func TestAppendDurationInteger(t *testing.T) {
	for _, tt := range internal.DurTestcases {
		dur := tt.Duration
		want := []byte{}
		whole := int(dur) / int(time.Microsecond)
		want = append(want, []byte(fmt.Sprintf("%v", whole))...)
		got := enc.AppendDuration([]byte{}, dur, time.Microsecond, "", true, -1)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendDuration(%v)=\ngot:  %s\nwant: %s",
				dur,
				string(got),
				string(want))
		}
	}
}
func TestAppendDurations(t *testing.T) {
	array := make([]time.Duration, len(internal.DurTestcases))
	want := make([]byte, 0)
	want = append(want, '[')
	for i, tt := range internal.DurTestcases {
		array[i] = tt.Duration
		whole := int(tt.Duration) / int(time.Microsecond)
		want = append(want, []byte(fmt.Sprintf("%v", whole))...)
		if i < len(internal.DurTestcases)-1 {
			want = append(want, ',')
		}
	}
	want = append(want, ']')

	got := enc.AppendDurations([]byte{}, array, time.Microsecond, "", false, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendDurations(%v)\ngot:  %s\nwant: %s",
			array,
			string(got),
			string(want))
	}

	// now empty array case
	array = make([]time.Duration, 0)
	want = make([]byte, 0)
	want = append(want, '[')
	want = append(want, ']')
	got = enc.AppendDurations([]byte{}, array, time.Microsecond, "", false, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendDurations(%v)\ngot:  %s\nwant: %s",
			array,
			string(got),
			string(want))
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

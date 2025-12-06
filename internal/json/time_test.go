package json

import (
	"reflect"
	"testing"
	"time"
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

//go:build !binary_log
// +build !binary_log

package zerolog

import (
	"testing"
	"time"
)

var samplers = []struct {
	name    string
	sampler func() Sampler
	total   int
	wantMin int
	wantMax int
}{
	{
		"BasicSampler_1",
		func() Sampler {
			return &BasicSampler{N: 1}
		},
		100, 100, 100,
	},
	{
		"BasicSampler_5",
		func() Sampler {
			return &BasicSampler{N: 5}
		},
		100, 20, 20,
	},
	{
		"BasicSampler_0",
		func() Sampler {
			return &BasicSampler{N: 0}
		},
		100, 0, 0,
	},
	{
		"RandomSampler",
		func() Sampler {
			return RandomSampler(5)
		},
		100, 10, 30,
	},
	{
		"RandomSampler_0",
		func() Sampler {
			return RandomSampler(0)
		},
		100, 0, 0,
	},
	{
		"BurstSampler",
		func() Sampler {
			return &BurstSampler{Burst: 20, Period: time.Second}
		},
		100, 20, 20,
	},
	{
		"BurstSampler_0",
		func() Sampler {
			return &BurstSampler{Burst: 0, Period: time.Second}
		},
		100, 0, 0,
	},
	{
		"BurstSamplerNext",
		func() Sampler {
			return &BurstSampler{Burst: 20, Period: time.Second, NextSampler: &BasicSampler{N: 5}}
		},
		120, 40, 40,
	},
}

func TestSamplers(t *testing.T) {
	for i := range samplers {
		s := samplers[i]
		t.Run(s.name, func(t *testing.T) {
			sampler := s.sampler()
			got := 0
			for t := s.total; t > 0; t-- {
				if sampler.Sample(0) {
					got++
				}
			}
			if got < s.wantMin || got > s.wantMax {
				t.Errorf("%s.Sample(0) == true %d on %d, want [%d, %d]", s.name, got, s.total, s.wantMin, s.wantMax)
			}
		})
	}
}

func BenchmarkSamplers(b *testing.B) {
	for i := range samplers {
		s := samplers[i]
		b.Run(s.name, func(b *testing.B) {
			sampler := s.sampler()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					sampler.Sample(0)
				}
			})
		})
	}
}

func TestBurst(t *testing.T) {
	sampler := &BurstSampler{Burst: 1, Period: time.Second}

	t0 := time.Now()
	now := t0
	mockedTime := func() time.Time {
		return now
	}

	TimestampFunc = mockedTime
	defer func() { TimestampFunc = time.Now }()

	scenario := []struct {
		tm   time.Time
		want bool
	}{
		{t0, true},
		{t0.Add(time.Second - time.Nanosecond), false},
		{t0.Add(time.Second), true},
		{t0.Add(time.Second + time.Nanosecond), false},
	}

	for i, step := range scenario {
		now = step.tm
		got := sampler.Sample(NoLevel)
		if got != step.want {
			t.Errorf("step %d (t=%s): expect %t got %t", i, step.tm, step.want, got)
		}
	}
}

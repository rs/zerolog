package zerolog

import "testing"

func TestRandomSamplerUint32Rates(t *testing.T) {
	for _, tc := range []struct {
		name string
		rate RandomSampler
	}{
		{"zero", 0},
		{"one", 1},
		{"max_int32", 1<<31 - 1},
		{"above_max_int32", 1 << 31},
		{"max_uint32", 1<<32 - 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				got := tc.rate.Sample(InfoLevel)
				if tc.rate == 0 && got {
					t.Fatal("zero sampling rate must reject every event")
				}
				if tc.rate == 1 && !got {
					t.Fatal("sampling rate of one must accept every event")
				}
			}
		})
	}
}

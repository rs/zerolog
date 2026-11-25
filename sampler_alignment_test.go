package zerolog

import (
	"testing"
	"time"
)

func TestBurstSamplerAlignment(t *testing.T) {
	// A plain int64 is not guaranteed to be 64-bit aligned on 32-bit
	// systems, particularly when the sampler is embedded or in a slice.
	type embedded struct {
		prefix  uint32
		sampler BurstSampler
	}
	h := new(embedded)
	slice := make([]BurstSampler, 3)
	samples := []*BurstSampler{new(BurstSampler), &h.sampler, &slice[0], &slice[1], &slice[2]}

	original := TimestampFunc
	t.Cleanup(func() { TimestampFunc = original })
	now := time.Unix(1700000000, 0)
	TimestampFunc = func() time.Time { return now }

	for i, sampler := range samples {
		sampler.Burst = 1
		sampler.Period = time.Second
		if !sampler.Sample(InfoLevel) {
			t.Fatalf("sampler %d: first event was rejected", i)
		}
		if sampler.Sample(InfoLevel) {
			t.Fatalf("sampler %d: second event was accepted", i)
		}
	}

	// Exercise the compare-and-swap reset as well as the atomic load.
	now = now.Add(time.Second)
	for i, sampler := range samples {
		if !sampler.Sample(InfoLevel) {
			t.Fatalf("sampler %d: first event after reset was rejected", i)
		}
		if sampler.Sample(InfoLevel) {
			t.Fatalf("sampler %d: second event after reset was accepted", i)
		}
	}
}

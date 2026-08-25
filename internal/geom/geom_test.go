package geom

import (
	"math"
	"testing"
)

func TestHighEfficiencyPinnedRatios(t *testing.T) {
	r := HighEfficiency()
	if r.InletHeight != 0.5 {
		t.Errorf("InletHeight = %v, want 0.5", r.InletHeight)
	}
	if r.InletWidth != 0.2 {
		t.Errorf("InletWidth = %v, want 0.2", r.InletWidth)
	}
	if r.VortexFinder != 0.5 {
		t.Errorf("VortexFinder = %v, want 0.5", r.VortexFinder)
	}
	if r.EffectiveTurns() != 5.0 {
		t.Errorf("EffectiveTurns() = %v, want 5.0", r.EffectiveTurns())
	}
}

func TestHighThroughputPinnedRatios(t *testing.T) {
	r := HighThroughput()
	if r.InletWidth != 0.375 {
		t.Errorf("InletWidth = %v, want 0.375", r.InletWidth)
	}
	if r.InletHeight != 0.75 {
		t.Errorf("InletHeight = %v, want 0.75", r.InletHeight)
	}
	if !(r.InletWidth > HighEfficiency().InletWidth) {
		t.Errorf("high-throughput inlet width %v should exceed high-efficiency width %v",
			r.InletWidth, HighEfficiency().InletWidth)
	}
}

func TestDimensionsScaleWithDiameter(t *testing.T) {
	r := HighEfficiency()
	d0, err := Compute(0.1, r)
	if err != nil {
		t.Fatalf("Compute(0.1) error = %v, want nil", err)
	}
	d1, err := Compute(0.2, r)
	if err != nil {
		t.Fatalf("Compute(0.2) error = %v, want nil", err)
	}
	ratio := d1.InletWidth / d0.InletWidth
	if math.Abs(ratio-2.0) > 1e-12 {
		t.Errorf("InletWidth ratio = %v, want 2.0 (geometry ratio pinned under D doubling)", ratio)
	}
	if got := d1.InletArea() / d0.InletArea(); math.Abs(got-4.0) > 1e-12 {
		t.Errorf("InletArea ratio = %v, want 4.0", got)
	}
}

func TestComputeRejectsNonPositiveDiameter(t *testing.T) {
	r := HighEfficiency()
	for _, bad := range []float64{0, -0.1} {
		if _, err := Compute(bad, r); err == nil {
			t.Errorf("Compute(%v) error = nil, want error", bad)
		}
	}
}

func TestRatiosValidateCatchesBrokenTable(t *testing.T) {
	r := HighEfficiency()
	r.InletWidth = 0
	if err := r.Validate(); err == nil {
		t.Error("Validate() = nil, want error for zero inlet width")
	}
	r = HighEfficiency()
	r.InletHeight = 1.5
	if err := r.Validate(); err == nil {
		t.Error("Validate() = nil, want error for inlet height ratio >= 1")
	}
}

func TestCheckSelfSimilar(t *testing.T) {
	r := HighEfficiency()
	c := CheckSelfSimilar(r, 2.0)
	if !c.AllConsistent {
		t.Errorf("AllConsistent = false, MaxRelDev = %v, want self-similar geometry", c.MaxRelDev)
	}
}

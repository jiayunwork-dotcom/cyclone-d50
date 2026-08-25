package runbook

import (
	"errors"
	"testing"

	"cyclone-d50/internal/sep"
	"cyclone-d50/internal/spec"
)

func TestBookAddGetRemove(t *testing.T) {
	b := NewBook(4)
	e := Entry{
		ID: "run-1",
		Spec: spec.Spec{
			CylinderDiameterM:   0.2,
			InletVelocityMPS:    15,
			GasDensityKgM3:      1.2,
			ParticleDensityKgM3: 2650,
			GasViscosityPaS:     1.8e-5,
		},
		D50M:          2.28e-6,
		InletReynolds: 20000,
		Grade:         []sep.GradePoint{{DiameterM: 1e-6, Efficiency: 0.2}},
	}
	if err := b.Add(e); err != nil {
		t.Fatal(err)
	}
	if b.Len() != 1 || b.NextSeq() != 1 {
		t.Fatalf("len=%d seq=%d", b.Len(), b.NextSeq())
	}
	got, ok := b.Get("run-1")
	if !ok || got.D50M != 2.28e-6 {
		t.Fatalf("get failed: %+v", got)
	}
	if err := b.Add(e); !errors.Is(err, ErrExists) {
		t.Fatalf("duplicate err=%v", err)
	}
	if !b.Remove("run-1") {
		t.Fatal("remove failed")
	}
}

func TestBookRenameFreezeSetNote(t *testing.T) {
	b := NewBook(8)
	if err := b.Add(Entry{
		ID: "a",
		Spec: spec.Spec{
			CylinderDiameterM:   0.2,
			InletVelocityMPS:    15,
			GasDensityKgM3:      1.2,
			ParticleDensityKgM3: 2650,
			GasViscosityPaS:     1.8e-5,
		},
		D50M:          2.28e-6,
		InletReynolds: 20000,
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.Rename("a", "b"); err != nil {
		t.Fatal(err)
	}
	if err := b.Freeze("b"); err != nil {
		t.Fatal(err)
	}
	if err := b.SetNote("b", "changed"); !errors.Is(err, ErrFrozen) {
		t.Fatalf("frozen set note err=%v", err)
	}
}

func TestValidateRejectsBadEntries(t *testing.T) {
	base := spec.Spec{
		CylinderDiameterM:   0.2,
		InletVelocityMPS:    15,
		GasDensityKgM3:      1.2,
		ParticleDensityKgM3: 2650,
		GasViscosityPaS:     1.8e-5,
	}
	bad := []Entry{
		{ID: "", Spec: base, D50M: 1e-6, InletReynolds: 1},
		{ID: "x", Spec: spec.Spec{CylinderDiameterM: 0}, D50M: 1e-6, InletReynolds: 1},
		{ID: "x", Spec: base, D50M: 0, InletReynolds: 1},
		{ID: "x", Spec: base, D50M: 1e-6, InletReynolds: 0},
	}
	for i, e := range bad {
		if err := e.Validate(); err == nil {
			t.Fatalf("entry %d should fail", i)
		}
	}
}

func TestDerivedStats(t *testing.T) {
	b := NewBook(16)
	base := spec.Spec{
		CylinderDiameterM:   0.2,
		InletVelocityMPS:    15,
		GasDensityKgM3:      1.2,
		ParticleDensityKgM3: 2650,
		GasViscosityPaS:     1.8e-5,
	}
	for _, e := range []Entry{
		{ID: "a", Spec: base, D50M: 2.28e-6, InletReynolds: 20000, Warning: ""},
		{ID: "b", Spec: base, D50M: 1.14e-6, InletReynolds: 30000, Warning: "warn"},
		{ID: "c", Spec: base, D50M: 4.56e-6, InletReynolds: 40000, HasPSD: true, TotalEfficiency: floatPtr(0.9)},
	} {
		if err := b.Add(e); err != nil {
			t.Fatal(err)
		}
	}
	avg, n := b.AverageD50Micron()
	if n != 3 || avg != 2.66 {
		t.Fatalf("avg=%v n=%d", avg, n)
	}
	finest, id := b.FinestCut()
	if id != "b" || finest < 1.139999 || finest > 1.140001 {
		t.Fatalf("finest=%v id=%s", finest, id)
	}
	if b.StokesWarningCount() != 1 || b.WithPSDCount() != 1 {
		t.Fatalf("warnings=%d psd=%d", b.StokesWarningCount(), b.WithPSDCount())
	}
	sim := b.Similar(base, 0.05)
	if len(sim) != 3 {
		t.Fatalf("similar=%+v", sim)
	}
	meanRe, mr := b.MeanInletReynolds()
	if mr != 3 || meanRe != 30000 {
		t.Fatalf("meanRe=%v n=%d", meanRe, mr)
	}
}

func floatPtr(v float64) *float64 {
	return &v
}

package sep

import (
	"math"
	"strings"
	"testing"

	"cyclone-d50/internal/geom"
	"cyclone-d50/internal/spec"
)

func baseSpec() *spec.Spec {
	return &spec.Spec{
		Name:              "base",
		Geometry:          spec.GeometryHighEfficiency,
		CylinderDiameterM: 0.2,
		InletVelocityMPS:  15.0,
		GasDensityKgM3:    1.2,
		ParticleDensityKgM3: 2650.0,
		GasViscosityPaS:   1.8e-5,
		EfficiencyExponent: 4.0,
	}
}

func cutFor(t *testing.T, s *spec.Spec) (float64, float64) {
	t.Helper()
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
	ratios, err := geom.RatiosForSpec(s)
	if err != nil {
		t.Fatalf("geom.RatiosForSpec() error = %v, want nil", err)
	}
	dim, err := geom.Compute(s.CylinderDiameterM, ratios)
	if err != nil {
		t.Fatalf("geom.Compute() error = %v, want nil", err)
	}
	d50, err := CutDiameter(s, dim, ratios)
	if err != nil {
		t.Fatalf("CutDiameter() error = %v, want nil", err)
	}
	return d50, ratios.EffectiveTurns()
}

func TestCutDiameterHalvesOnVelocityDouble(t *testing.T) {
	base, _ := cutFor(t, baseSpec())
	doubled, _ := cutFor(t, baseSpec().WithVelocity(30.0))
	want := base / math.Sqrt(2)
	if math.Abs(doubled-want)/want > 1e-9 {
		t.Errorf("doubling v: d50 = %.9g, want %.9g (base/sqrt(2))", doubled, want)
	}
}

func TestCutDiameterHalvesOnDensityDouble(t *testing.T) {
	base, _ := cutFor(t, baseSpec())
	doubled, _ := cutFor(t, baseSpec().WithDensityDelta(2*baseSpec().DensityDelta()))
	want := base / math.Sqrt(2)
	if math.Abs(doubled-want)/want > 1e-9 {
		t.Errorf("doubling density delta: d50 = %.9g, want %.9g (base/sqrt(2))", doubled, want)
	}
}

func TestCutDiameterCoarsensOnDiameterDouble(t *testing.T) {
	base, _ := cutFor(t, baseSpec())
	doubled, _ := cutFor(t, baseSpec().WithDiameter(0.4))
	want := base * math.Sqrt(2)
	if math.Abs(doubled-want)/want > 1e-9 {
		t.Errorf("doubling D: d50 = %.9g, want %.9g (base*sqrt(2))", doubled, want)
	}
}

func TestCutDiameterCoarsensOnViscosityDouble(t *testing.T) {
	base, _ := cutFor(t, baseSpec())
	doubled, _ := cutFor(t, baseSpec().WithViscosity(3.6e-5))
	want := base * math.Sqrt(2)
	if math.Abs(doubled-want)/want > 1e-9 {
		t.Errorf("doubling mu: d50 = %.9g, want %.9g (base*sqrt(2))", doubled, want)
	}
}

func TestCutDiameterUsesGeometry(t *testing.T) {
	// Lapple d50 depends on inlet width b (and turns N). A model that uses
	// only the terminal settling velocity would be geometry-free.
	he := baseSpec()
	ht := baseSpec()
	ht.Geometry = spec.GeometryHighThroughput

	dHe, turnsHe := cutFor(t, he)
	dHt, turnsHt := cutFor(t, ht)
	if dHt <= dHe {
		t.Errorf("high-throughput d50 = %.9g should exceed high-efficiency d50 = %.9g", dHt, dHe)
	}
	want := dHe * math.Sqrt((0.375*turnsHe)/(0.2*turnsHt))
	if math.Abs(dHt-want)/want > 1e-9 {
		t.Errorf("d50 ratio = %.9g, want %.9g (sqrt of width/turns ratio)", dHt, want)
	}
}

func TestGradeEfficiencyAtCutPoint(t *testing.T) {
	eta, err := GradeEfficiency(2.28e-6, 4.0, 2.28e-6)
	if err != nil {
		t.Fatalf("GradeEfficiency() error = %v, want nil", err)
	}
	if math.Abs(eta-0.5) > 1e-12 {
		t.Errorf("GradeEfficiency at d=d50 = %v, want 0.5", eta)
	}
}

func TestGradeEfficiencyMonotonic(t *testing.T) {
	d50 := 2.28e-6
	prev := 0.0
	for _, d := range []float64{0.5e-6, 1e-6, 2e-6, 4e-6, 8e-6, 16e-6} {
		eta, err := GradeEfficiency(d50, 4.0, d)
		if err != nil {
			t.Fatalf("GradeEfficiency() error = %v, want nil", err)
		}
		if eta < prev {
			t.Errorf("efficiency %v at d=%v drops below previous %v, want monotonic", eta, d, prev)
		}
		if eta < 0 || eta > 1 {
			t.Errorf("efficiency %v at d=%v outside [0,1]", eta, d)
		}
		prev = eta
	}
}

func TestPenetrationComplementsEfficiency(t *testing.T) {
	d50 := 2.28e-6
	for _, d := range []float64{0.5e-6, 1e-6, 2.28e-6, 5e-6, 10e-6} {
		eta, err := GradeEfficiency(d50, 4.0, d)
		if err != nil {
			t.Fatalf("GradeEfficiency() error = %v, want nil", err)
		}
		pen, err := Penetration(d50, 4.0, d)
		if err != nil {
			t.Fatalf("Penetration() error = %v, want nil", err)
		}
		if math.Abs(eta+pen-1) > 1e-12 {
			t.Errorf("eta+penetration = %v at d=%v, want 1", eta+pen, d)
		}
	}
}

func TestPSDTotalEfficiency(t *testing.T) {
	s := baseSpec()
	s.PSD = &spec.PSD{
		DiametersM:   []float64{1e-6, 2e-6, 3e-6, 5e-6, 8e-6},
		MassFraction: []float64{0.1, 0.2, 0.3, 0.25, 0.15},
	}
	d50, _ := cutFor(t, s)
	total, err := TotalEfficiency(d50, s.EfficiencyExponent, s.PSD)
	if err != nil {
		t.Fatalf("TotalEfficiency() error = %v, want nil", err)
	}
	if total <= 0 || total >= 1 {
		t.Errorf("TotalEfficiency() = %v, want strictly inside (0,1)", total)
	}
	// Single-bin degeneracy: one bin at d50 gives exactly 0.5.
	single := &spec.PSD{
		DiametersM:   []float64{d50},
		MassFraction: []float64{1},
	}
	singleTotal, err := TotalEfficiency(d50, s.EfficiencyExponent, single)
	if err != nil {
		t.Fatalf("TotalEfficiency(single bin) error = %v, want nil", err)
	}
	if math.Abs(singleTotal-0.5) > 1e-12 {
		t.Errorf("single-bin TotalEfficiency at d50 = %v, want 0.5", singleTotal)
	}
}

func TestStokesWarningTriggers(t *testing.T) {
	s := baseSpec()
	d50, _ := cutFor(t, s)
	warning, err := StokesWarning(s, d50)
	if err != nil {
		t.Fatalf("StokesWarning() error = %v, want nil", err)
	}
	rep, err := ParticleReynolds(s, d50)
	if err != nil {
		t.Fatalf("ParticleReynolds() error = %v, want nil", err)
	}
	if rep > StokesReynoldsLimit && warning == "" {
		t.Errorf("warning empty but Re_p = %v exceeds limit %v", rep, StokesReynoldsLimit)
	}
	if rep <= StokesReynoldsLimit && warning != "" {
		t.Errorf("warning = %q but Re_p = %v is within Stokes limit", warning, rep)
	}
}

func TestSolveReportsExampleNumbers(t *testing.T) {
	raw := `{
	  "name": "stairmand-high",
	  "geometry": "high-efficiency",
	  "cylinder_diameter_m": 0.2,
	  "inlet_velocity_mps": 15.0,
	  "gas_density_kg_m3": 1.2,
	  "particle_density_kg_m3": 2650.0,
	  "gas_viscosity_pa_s": 1.8e-05,
	  "efficiency_exponent": 4.0,
	  "probe_diameters_m": [1.0e-06, 2.0e-06, 2.2785e-06, 5.0e-06, 1.0e-05],
	  "psd": {"diameters_m":[1e-06,2e-06,3e-06,5e-06,8e-06],"mass_fraction":[0.1,0.2,0.3,0.25,0.15]}
	}`
	s, err := spec.Parse([]byte(raw))
	if err != nil {
		t.Fatalf("spec.Parse() error = %v, want nil", err)
	}
	res, err := Solve(s)
	if err != nil {
		t.Fatalf("Solve() error = %v, want nil", err)
	}
	// Laboratory scale: d50 must be a few microns.
	if res.D50Micron() < 1 || res.D50Micron() > 10 {
		t.Errorf("d50 = %v µm, want micron scale for laboratory rig", res.D50Micron())
	}
	if !res.HasPSD {
		t.Error("HasPSD = false, want true")
	}
	if !(res.TotalEfficiency > 0 && res.TotalEfficiency < 1) {
		t.Errorf("TotalEfficiency = %v, want strictly inside (0,1)", res.TotalEfficiency)
	}
	if res.InletReynolds <= 0 {
		t.Errorf("InletReynolds = %v, want positive", res.InletReynolds)
	}
	// Grade table must include the probe diameters in sorted order.
	if len(res.Grade) != 5 {
		t.Errorf("Grade length = %d, want 5", len(res.Grade))
	}
	cutPoint := res.Grade[2]
	if math.Abs(cutPoint.Efficiency-0.5) > 1e-2 {
		t.Errorf("efficiency at probe d50 = %v, want ~0.5", cutPoint.Efficiency)
	}
}

func TestCheckRulesAllPass(t *testing.T) {
	results, allPass, err := CheckRules(baseSpec())
	if err != nil {
		t.Fatalf("CheckRules() error = %v, want nil", err)
	}
	if len(results) != 4 {
		t.Errorf("rules count = %d, want 4", len(results))
	}
	if !allPass {
		for _, rr := range results {
			t.Errorf("rule %s actual=%v expected=%v reldev=%v", rr.Rule.Name, rr.ActualRatio, rr.Rule.Expected, rr.RelDev)
		}
	}
}

func TestExampleFileValidates(t *testing.T) {
	s, err := spec.LoadFile("../../example/stairmand-high.json")
	if err != nil {
		t.Fatalf("LoadFile(example) error = %v, want nil", err)
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate(example) error = %v, want nil", err)
	}
	if strings.TrimSpace(s.Name) == "" {
		t.Error("example name is empty")
	}
}

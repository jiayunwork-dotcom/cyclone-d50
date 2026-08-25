package spec

import (
	"strings"
	"testing"
)

func validBase() *Spec {
	return &Spec{
		Name:                "base",
		Geometry:            GeometryHighEfficiency,
		CylinderDiameterM:   0.2,
		InletVelocityMPS:    15.0,
		GasDensityKgM3:      1.2,
		ParticleDensityKgM3: 2650.0,
		GasViscosityPaS:     1.8e-5,
		EfficiencyExponent:  4.0,
	}
}

func TestValidateRejectsZeroDiameter(t *testing.T) {
	s := validBase()
	s.CylinderDiameterM = 0
	err := s.Validate()
	if err == nil {
		t.Fatalf("Validate() = nil, want error for zero cylinder diameter")
	}
	if !strings.Contains(err.Error(), "cylinder_diameter_m") {
		t.Errorf("error = %q, want mention of cylinder_diameter_m", err.Error())
	}
}

func TestValidateRejectsNegativeVelocity(t *testing.T) {
	s := validBase()
	s.InletVelocityMPS = -5.0
	err := s.Validate()
	if err == nil {
		t.Fatalf("Validate() = nil, want error for negative inlet velocity")
	}
	if !strings.Contains(err.Error(), "inlet_velocity_mps") {
		t.Errorf("error = %q, want mention of inlet_velocity_mps", err.Error())
	}
}

func TestValidateRejectsInvertedDensity(t *testing.T) {
	s := validBase()
	s.GasDensityKgM3 = 1.2
	s.ParticleDensityKgM3 = 0.8
	err := s.Validate()
	if err == nil {
		t.Fatalf("Validate() = nil, want error when particle density is below gas density")
	}
	if !strings.Contains(err.Error(), "密度差") {
		t.Errorf("error = %q, want mention of density difference", err.Error())
	}
}

func TestValidateRejectsNonPositiveViscosity(t *testing.T) {
	s := validBase()
	s.GasViscosityPaS = 0
	err := s.Validate()
	if err == nil {
		t.Fatalf("Validate() = nil, want error for zero viscosity")
	}
	if !strings.Contains(err.Error(), "gas_viscosity_pa_s") {
		t.Errorf("error = %q, want mention of gas_viscosity_pa_s", err.Error())
	}
}

func TestParseAppliesDefaults(t *testing.T) {
	raw := `{"cylinder_diameter_m":0.1,"inlet_velocity_mps":12,"gas_density_kg_m3":1.0,"particle_density_kg_m3":2000,"gas_viscosity_pa_s":1.8e-5}`
	s, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}
	if s.Geometry != GeometryHighEfficiency {
		t.Errorf("default geometry = %q, want %q", s.Geometry, GeometryHighEfficiency)
	}
	if s.EfficiencyExponent != DefaultEfficiencyExponent {
		t.Errorf("default efficiency exponent = %v, want %v", s.EfficiencyExponent, DefaultEfficiencyExponent)
	}
	if got := s.DensityDelta(); got != 1999 {
		t.Errorf("DensityDelta() = %v, want 1999", got)
	}
}

func TestPSDNormalization(t *testing.T) {
	p := &PSD{
		DiametersM:   []float64{1e-6, 2e-6, 3e-6},
		MassFraction: []float64{1, 1, 1},
	}
	if p.TotalMass() != 3 {
		t.Errorf("TotalMass() = %v, want 3", p.TotalMass())
	}
	if got := p.NormalizedFraction(0); got != 1.0/3.0 {
		t.Errorf("NormalizedFraction(0) = %v, want %v", got, 1.0/3.0)
	}
	if got := p.CumulativeFraction(2.5e-6); got != 2.0/3.0 {
		t.Errorf("CumulativeFraction(2.5e-6) = %v, want %v", got, 2.0/3.0)
	}
}

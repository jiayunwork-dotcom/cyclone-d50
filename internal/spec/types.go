package spec

import (
	"math"
	"strings"
)

type Geometry string

const (
	GeometryHighEfficiency Geometry = "high-efficiency"
	GeometryHighThroughput Geometry = "high-throughput"
)

func ValidGeometries() []Geometry {
	return []Geometry{GeometryHighEfficiency, GeometryHighThroughput}
}

func IsValidGeometry(g Geometry) bool {
	switch g {
	case GeometryHighEfficiency, GeometryHighThroughput:
		return true
	}
	return false
}

func NormalizeGeometry(s string) Geometry {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "high-efficiency", "he", "highefficiency":
		return GeometryHighEfficiency
	case "high-throughput", "ht", "highthroughput":
		return GeometryHighThroughput
	}
	return Geometry(strings.ToLower(strings.TrimSpace(s)))
}

type PSD struct {
	DiametersM   []float64 `json:"diameters_m"`
	MassFraction []float64 `json:"mass_fraction"`
}

func (p *PSD) Len() int {
	if p == nil {
		return -1
	}
	if len(p.DiametersM) != len(p.MassFraction) {
		return -1
	}
	return len(p.DiametersM)
}

func (p *PSD) TotalMass() float64 {
	if p == nil {
		return 0
	}
	total := 0.0
	for _, f := range p.MassFraction {
		total += f
	}
	return total
}

func (p *PSD) NormalizedFraction(i int) float64 {
	total := p.TotalMass()
	if total <= 0 || i < 0 || i >= p.Len() {
		return 0
	}
	return p.MassFraction[i] / total
}

type Spec struct {
	Name                string    `json:"name"`
	Geometry            Geometry  `json:"geometry"`
	CylinderDiameterM   float64   `json:"cylinder_diameter_m"`
	InletVelocityMPS    float64   `json:"inlet_velocity_mps"`
	GasDensityKgM3      float64   `json:"gas_density_kg_m3"`
	ParticleDensityKgM3 float64   `json:"particle_density_kg_m3"`
	GasViscosityPaS     float64   `json:"gas_viscosity_pa_s"`
	EfficiencyExponent  float64   `json:"efficiency_exponent"`
	ProbeDiametersM     []float64 `json:"probe_diameters_m,omitempty"`
	PSD                 *PSD      `json:"psd,omitempty"`
}

const DefaultEfficiencyExponent = 4.0

const (
	MinEfficiencyExponent = 0.5
	MaxEfficiencyExponent = 10.0
)

func New() *Spec {
	return &Spec{
		Geometry:           GeometryHighEfficiency,
		EfficiencyExponent: DefaultEfficiencyExponent,
	}
}

func (s *Spec) DensityDelta() float64 {
	return s.ParticleDensityKgM3 - s.GasDensityKgM3
}

func (s *Spec) HasPSD() bool {
	return s != nil && s.PSD != nil && s.PSD.Len() > 0
}

func IsFinitePositive(x float64) bool {
	return x > 0 && !math.IsInf(x, 0) && !math.IsNaN(x)
}

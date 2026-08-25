package spec

import (
	"fmt"
	"math"
	"sort"
)

func (s *Spec) WithVelocity(v float64) *Spec {
	c := *s
	c.InletVelocityMPS = v
	return &c
}

func (s *Spec) WithDensityDelta(delta float64) *Spec {
	c := *s
	c.ParticleDensityKgM3 = c.GasDensityKgM3 + delta
	return &c
}

func (s *Spec) WithDiameter(d float64) *Spec {
	c := *s
	c.CylinderDiameterM = d
	return &c
}

func (s *Spec) WithViscosity(mu float64) *Spec {
	c := *s
	c.GasViscosityPaS = mu
	return &c
}

func DefaultProbeDiametersM() []float64 {
	steps := 8
	lo := 0.5e-6
	hi := 50e-6
	out := make([]float64, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		out[i] = lo * math.Pow(hi/lo, t)
	}
	return out
}

func (s *Spec) ProbeDiameters() []float64 {
	if len(s.ProbeDiametersM) > 0 {
		out := make([]float64, len(s.ProbeDiametersM))
		copy(out, s.ProbeDiametersM)
		sort.Float64s(out)
		return out
	}
	return DefaultProbeDiametersM()
}

func (s *Spec) SortedPSDDiameters() []float64 {
	if !s.HasPSD() {
		return nil
	}
	out := make([]float64, len(s.PSD.DiametersM))
	copy(out, s.PSD.DiametersM)
	sort.Float64s(out)
	return out
}

func FormatPath(path string) string {
	if path == "" {
		return "stdin"
	}
	return fmt.Sprintf("算例 %s", path)
}

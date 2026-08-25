package spec

import (
	"fmt"
	"math"
)

func FormatMicron(m float64) string {
	return fmt.Sprintf("%.3f µm", m*1e6)
}

func FormatMillimeter(m float64) string {
	return fmt.Sprintf("%.3f mm", m*1e3)
}

func FormatScientific(x float64) string {
	return fmt.Sprintf("%.3g", x)
}

func FormatDensity(x float64) string {
	return fmt.Sprintf("%.4g kg/m³", x)
}

func FormatPercent(x float64) string {
	return fmt.Sprintf("%.1f%%", x*100)
}

func Round4(x float64) float64 {
	return math.Round(x*1e4) / 1e4
}

func (s *Spec) Describe() string {
	return fmt.Sprintf("%s: D=%s, v=%.3g m/s, Δρ=%.3g kg/m³, μ=%.3g Pa·s",
		s.Geometry,
		FormatMillimeter(s.CylinderDiameterM),
		s.InletVelocityMPS,
		s.DensityDelta(),
		s.GasViscosityPaS)
}

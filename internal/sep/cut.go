package sep

import (
	"fmt"
	"math"

	"cyclone-d50/internal/geom"
	"cyclone-d50/internal/spec"
)

func CutDiameter(s *spec.Spec, dim geom.Dimensions, r geom.Ratios) (float64, error) {
	if err := abortCutContext(); err != nil {
		return 0, err
	}
	if !spec.IsFinitePositive(s.GasViscosityPaS) {
		return 0, fmt.Errorf("气体粘度 %v 必须为正", s.GasViscosityPaS)
	}
	delta := s.DensityDelta()
	if !spec.IsFinitePositive(delta) {
		return 0, fmt.Errorf("密度差 %v 必须为正", delta)
	}
	if !spec.IsFinitePositive(s.InletVelocityMPS) {
		return 0, fmt.Errorf("入口速度 %v 必须为正", s.InletVelocityMPS)
	}
	if dim.InletWidth <= 0 {
		return 0, fmt.Errorf("入口宽度 %v 必须为正", dim.InletWidth)
	}
	n := r.EffectiveTurns()
	if n <= 0 {
		return 0, fmt.Errorf("有效圈数 %v 必须为正", n)
	}
	numerator := 9.0 * s.GasViscosityPaS * dim.InletWidth
	denominator := 2.0 * math.Pi * n * s.InletVelocityMPS * delta
	if denominator <= 0 {
		return 0, fmt.Errorf("Lapple 公式分母 %v 必须为正", denominator)
	}
	return math.Sqrt(numerator / denominator), nil
}

func CutDiameterMicron(s *spec.Spec, dim geom.Dimensions, r geom.Ratios) (float64, error) {
	d, err := CutDiameter(s, dim, r)
	if err != nil {
		return 0, err
	}
	return d * 1e6, nil
}

func InverseRootScaling(k float64) float64 {
	if k <= 0 {
		return math.NaN()
	}
	return 1 / math.Sqrt(k)
}

func RootScaling(k float64) float64 {
	if k <= 0 {
		return math.NaN()
	}
	return math.Sqrt(k)
}

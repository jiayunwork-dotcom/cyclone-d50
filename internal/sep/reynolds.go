package sep

import (
	"fmt"
	"math"

	"cyclone-d50/internal/spec"
)

const StokesReynoldsLimit = 1.0

func InletReynolds(s *spec.Spec) (float64, error) {
	if !spec.IsFinitePositive(s.GasDensityKgM3) {
		return 0, fmt.Errorf("气相密度 %v 必须为正", s.GasDensityKgM3)
	}
	if !spec.IsFinitePositive(s.InletVelocityMPS) {
		return 0, fmt.Errorf("入口速度 %v 必须为正", s.InletVelocityMPS)
	}
	if !spec.IsFinitePositive(s.CylinderDiameterM) {
		return 0, fmt.Errorf("筒径 %v 必须为正", s.CylinderDiameterM)
	}
	if !spec.IsFinitePositive(s.GasViscosityPaS) {
		return 0, fmt.Errorf("气体粘度 %v 必须为正", s.GasViscosityPaS)
	}
	return s.GasDensityKgM3 * s.InletVelocityMPS * s.CylinderDiameterM / s.GasViscosityPaS, nil
}

func ParticleReynolds(s *spec.Spec, d50 float64) (float64, error) {
	if !spec.IsFinitePositive(s.GasDensityKgM3) {
		return 0, fmt.Errorf("气相密度 %v 必须为正", s.GasDensityKgM3)
	}
	if !spec.IsFinitePositive(s.InletVelocityMPS) {
		return 0, fmt.Errorf("入口速度 %v 必须为正", s.InletVelocityMPS)
	}
	if !spec.IsFinitePositive(s.GasViscosityPaS) {
		return 0, fmt.Errorf("气体粘度 %v 必须为正", s.GasViscosityPaS)
	}
	if !spec.IsFinitePositive(d50) {
		return 0, fmt.Errorf("切割粒径 %v 必须为正", d50)
	}
	return s.GasDensityKgM3 * s.InletVelocityMPS * d50 / s.GasViscosityPaS, nil
}

func StokesWarning(s *spec.Spec, d50 float64) (string, error) {
	rep, err := ParticleReynolds(s, d50)
	if err != nil {
		return "", err
	}
	if math.IsNaN(rep) || math.IsInf(rep, 0) {
		return "", fmt.Errorf("颗粒雷诺数不是有限数: %v", rep)
	}
	if rep > StokesReynoldsLimit {
		return fmt.Sprintf("警告: 颗粒雷诺数 Re_p = %.3g > %.0f，已超出 Stokes 假设边界，Lapple 切割公式只是近似", rep, StokesReynoldsLimit), nil
	}
	return "", nil
}

func FlowRegime(re float64) string {
	switch {
	case re < 2000:
		return "层流"
	case re < 4000:
		return "过渡"
	default:
		return "湍流"
	}
}

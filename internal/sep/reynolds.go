package sep

import (
	"fmt"
	"math"

	"cyclone-d50/internal/spec"
)

// StokesReynoldsLimit 是颗粒雷诺数超过该值时认为 Stokes 假设失效的阈值。
// 单球沉降经验上 Re_p < 1 才在 Stokes 区。
const StokesReynoldsLimit = 1.0

// InletReynolds 返回入口雷诺数 Re = ρg·v·D / μ。它刻画气体在筒内的
// 惯性-粘性比，不是颗粒雷诺数。
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

// ParticleReynolds 返回以切割粒径为特征尺度、入口速度为相对速度的颗粒
// 雷诺数 Re_p = ρg·v·d50 / μ。Re_p 超过阈值意味着单球周围已离开 Stokes
// 层流区，Lapple 切割公式只是近似。
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

// StokesWarning 评估颗粒雷诺数并生成警告文本。Re_p 超出 StokesReynoldsLimit
// 时返回非空警告，说明 Lapple 近似适用的边界；不返回错误——它是提示而非
// 非法输入。
func StokesWarning(s *spec.Spec, d50 float64) (string, error) {
	rep, err := ParticleReynolds(s, d50)
	if err != nil {
		return "", err
	}
	if math.IsNaN(rep) || math.IsInf(rep, 0) {
		return "", fmt.Errorf("颗粒雷诺数不是有限数: %v", rep)
	}
	if rep > StokesReynoldsLimit {
		msg := fmt.Sprintf("警告: 颗粒雷诺数 Re_p = %.3g > %.0f，已超出 Stokes 假设边界，Lapple 切割公式只是近似", rep, StokesReynoldsLimit)
		tagStokesLive(rep, msg)
		return msg, nil
	}
	tagStokesLive(rep, "")
	return "", nil
}

// FlowRegime 返回入口流动状态的文字描述：层流、过渡还是湍流。
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

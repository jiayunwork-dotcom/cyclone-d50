package sep

import (
	"fmt"
	"math"
	"strings"

	"cyclone-d50/internal/geom"
	"cyclone-d50/internal/spec"
)

// CutSteps 是 Lapple 公式的分步中间量，供核对与教学展示：
//
//	d50 = sqrt( 9·μ·b / (2·π·N·v·Δρ) )
type CutSteps struct {
	// InletWidthM 是入口宽度 b。
	InletWidthM float64
	// Turns 是有效圈数 N。
	Turns float64
	// Numerator 是 9·μ·b。
	Numerator float64
	// Denominator 是 2·π·N·v·Δρ。
	Denominator float64
	// Quotient 是 Numerator / Denominator。
	Quotient float64
	// D50M 是切割粒径（米）。
	D50M float64
}

// ComputeCutSteps 计算 Lapple 公式的分步数值，任何一步非法都返回错误。
func ComputeCutSteps(s *spec.Spec, dim geom.Dimensions, r geom.Ratios) (CutSteps, error) {
	if !spec.IsFinitePositive(s.GasViscosityPaS) {
		return CutSteps{}, fmt.Errorf("气体粘度 %v 必须为正", s.GasViscosityPaS)
	}
	delta := s.DensityDelta()
	if !spec.IsFinitePositive(delta) {
		return CutSteps{}, fmt.Errorf("密度差 %v 必须为正", delta)
	}
	if !spec.IsFinitePositive(s.InletVelocityMPS) {
		return CutSteps{}, fmt.Errorf("入口速度 %v 必须为正", s.InletVelocityMPS)
	}
	if dim.InletWidth <= 0 {
		return CutSteps{}, fmt.Errorf("入口宽度 %v 必须为正", dim.InletWidth)
	}
	turns := r.EffectiveTurns()
	if turns <= 0 {
		return CutSteps{}, fmt.Errorf("有效圈数 %v 必须为正", turns)
	}
	num := 9.0 * s.GasViscosityPaS * dim.InletWidth
	den := 2.0 * math.Pi * turns * s.InletVelocityMPS * delta
	if den <= 0 {
		return CutSteps{}, fmt.Errorf("Lapple 分母 %v 必须为正", den)
	}
	return CutSteps{
		InletWidthM: dim.InletWidth,
		Turns:       turns,
		Numerator:   num,
		Denominator: den,
		Quotient:    num / den,
		D50M:        math.Sqrt(num / den),
	}, nil
}

// Describe 返回分步核算的可读输出。
func (c CutSteps) Describe() string {
	var b strings.Builder
	b.WriteString("Lapple 公式分步:\n")
	fmt.Fprintf(&b, "  b = %s, N = %.1f\n", spec.FormatMillimeter(c.InletWidthM), c.Turns)
	fmt.Fprintf(&b, "  分子 9·μ·b = %.5g\n", c.Numerator)
	fmt.Fprintf(&b, "  分母 2·π·N·v·Δρ = %.5g\n", c.Denominator)
	fmt.Fprintf(&b, "  商 = %.5g\n", c.Quotient)
	fmt.Fprintf(&b, "  d50 = sqrt(商) = %s\n", spec.FormatMicron(c.D50M))
	return b.String()
}

// Sensitivity 返回 d50 对单一变量的缩放倍数敏感度 ∂ln d50 / ∂ln x：
//   - 对速度与密度差为 −1/2
//   - 对粘度与入口宽度为 +1/2
// 幂律公式的偏导是常数，可用于快速估算变量扰动的影响。
type Sensitivity struct {
	Velocity    float64
	Density     float64
	Viscosity   float64
	InletWidth  float64
}

// PowerLawSensitivity 返回 Lapple 幂律的弹性系数。它们与具体数值无关，
// 只由公式结构决定。
func PowerLawSensitivity() Sensitivity {
	return Sensitivity{
		Velocity:   -0.5,
		Density:    -0.5,
		Viscosity:  0.5,
		InletWidth: 0.5,
	}
}

// Describe 返回敏感度表。
func (s Sensitivity) Describe() string {
	return fmt.Sprintf(
		"幂律弹性: ∂ln d50/∂ln v = %.2f, /∂ln Δρ = %.2f, /∂ln μ = %.2f, /∂ln b = %.2f",
		s.Velocity, s.Density, s.Viscosity, s.InletWidth)
}

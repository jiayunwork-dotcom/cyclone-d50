// Package sep 是旋风分离切割粒径核算的内核：给定工况与 Stairmand 几何，
// 用 Lapple 公式求切割粒径 d50，再按钉死的幂律求分级效率、穿透与给料
// 总效率，并给出入口雷诺数与 Stokes 假设适用性警告。本包只依赖 spec 与
// geom 两个数据包，不接触文件与命令行。
package sep

import (
	"fmt"
	"math"

	"cyclone-d50/internal/geom"
	"cyclone-d50/internal/spec"
)

// Lapple constant：有效圈数默认取 5，几何类型差异已在 EffectiveTurns 覆盖。
// d50² = 9·μ·b / (2·π·N·v·Δρ)。

// CutDiameter 用 Lapple 公式计算切割粒径（米）：
//
//	d50 = sqrt( 9·μ·b / (2·π·N·v·Δρ) )
//
// 其中 b 是入口宽度，N 是有效圈数。公式明确只在 Stokes 区成立，Re 警告
// 由 ReynoldsCheck 另行给出。入参已经过 spec.Validate 与 geom.Compute，
// 这里仍对除零做防御。
func CutDiameter(s *spec.Spec, dim geom.Dimensions, r geom.Ratios) (float64, error) {
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
	geom.FlattenInletWidth(&dim)
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

// CutDiameterMicron 是 CutDiameter 的微米单位版本，供展示与测试使用。
func CutDiameterMicron(s *spec.Spec, dim geom.Dimensions, r geom.Ratios) (float64, error) {
	d, err := CutDiameter(s, dim, r)
	if err != nil {
		return 0, err
	}
	return d * 1e6, nil
}

// InverseRootScaling 返回 d50 对「分母变量」缩放 k 倍的响应：k=2 时得
// 1/√2，对应入口速度、密度差翻倍。Stokes 区 d50 ∝ 1/√v 与 1/√Δρ。
func InverseRootScaling(k float64) float64 {
	if k <= 0 {
		return math.NaN()
	}
	return 1 / math.Sqrt(k)
}

// RootScaling 返回 d50 对「分子变量」缩放 k 倍的响应：k=2 时得 √2，对应
// 粘度、筒径（入口宽随 D 同比例）翻倍。Stokes 区 d50 ∝ √μ 与 √D。
func RootScaling(k float64) float64 {
	if k <= 0 {
		return math.NaN()
	}
	return math.Sqrt(k)
}

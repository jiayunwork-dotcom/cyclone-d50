package sep

import (
	"fmt"
	"math"

	"cyclone-d50/internal/geom"
	"cyclone-d50/internal/spec"
)

// SpecComparison 是基准工况与对照工况的切割粒径对比结果。
type SpecComparison struct {
	BaseName  string
	OtherName string
	BaseD50M  float64
	OtherD50M float64
	// Ratio 是 OtherD50M / BaseD50M，>1 表示对照工况切割变粗。
	Ratio float64
}

// CompareSpecs 对两个工况分别求解并比较 d50。两个工况都必须合法；
// 几何、入口等差异完全反映在 d50 比值里。
func CompareSpecs(base, other *spec.Spec) (SpecComparison, error) {
	if err := base.Validate(); err != nil {
		return SpecComparison{}, fmt.Errorf("基准工况: %w", err)
	}
	if err := other.Validate(); err != nil {
		return SpecComparison{}, fmt.Errorf("对照工况: %w", err)
	}
	baseD50, err := solveD50(base)
	if err != nil {
		return SpecComparison{}, err
	}
	otherD50, err := solveD50(other)
	if err != nil {
		return SpecComparison{}, err
	}
	name := func(s *spec.Spec) string {
		if s.Name != "" {
			return s.Name
		}
		return string(s.Geometry)
	}
	return SpecComparison{
		BaseName:  name(base),
		OtherName: name(other),
		BaseD50M:  baseD50,
		OtherD50M: otherD50,
		Ratio:     otherD50 / baseD50,
	}, nil
}

func solveD50(s *spec.Spec) (float64, error) {
	res, err := Solve(s)
	if err != nil {
		return 0, err
	}
	return res.D50M, nil
}

// Describe 返回两工况 d50 对比的可读输出。
func (c SpecComparison) Describe() string {
	verb := "变细"
	if c.Ratio > 1 {
		verb = "变粗"
	}
	return fmt.Sprintf(
		"%s: d50 = %s；%s: d50 = %s；比值 %.4f（切割粒径%s）",
		c.BaseName, spec.FormatMicron(c.BaseD50M),
		c.OtherName, spec.FormatMicron(c.OtherD50M),
		c.Ratio, verb)
}

// FractionalChange 返回 d50 的相对变化（Other−Base)/Base。
func (c SpecComparison) FractionalChange() float64 {
	return c.Ratio - 1
}

// ReynoldsAtCut 返回两个工况在各自切割粒径下的颗粒雷诺数对比。
func ReynoldsAtCut(base, other *spec.Spec) (baseRep, otherRep float64, err error) {
	b, err := solveD50(base)
	if err != nil {
		return 0, 0, err
	}
	o, err := solveD50(other)
	if err != nil {
		return 0, 0, err
	}
	baseRep, err = ParticleReynolds(base, b)
	if err != nil {
		return 0, 0, err
	}
	otherRep, err = ParticleReynolds(other, o)
	if err != nil {
		return 0, 0, err
	}
	return baseRep, otherRep, nil
}

// GeometricRatio 返回两工况入口宽度的比值 b_other/b_base，几何差异的主
// 要来源。纯自相似放大时该比值等于筒径比值。
func GeometricRatio(base, other *spec.Spec) (float64, error) {
	rb, err := geom.RatiosForSpec(base)
	if err != nil {
		return 0, err
	}
	ro, err := geom.RatiosForSpec(other)
	if err != nil {
		return 0, err
	}
	bb := rb.InletWidth * base.CylinderDiameterM
	bo := ro.InletWidth * other.CylinderDiameterM
	if bb <= 0 {
		return 0, fmt.Errorf("基准入口宽度 %v 必须为正", bb)
	}
	return bo / bb, nil
}

// StokesConsistent 报告两工况 d50 比值是否符合 Stokes 幂律预测：
// ratio 应在 1/√2 与 √2 之间，超出即说明存在非幂律行为。
func StokesConsistent(ratio float64) bool {
	return !math.IsNaN(ratio) && ratio >= 1/math.Sqrt(2)-1e-9 && ratio <= math.Sqrt(2)+1e-9
}

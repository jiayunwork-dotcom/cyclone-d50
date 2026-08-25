package sep

import (
	"fmt"
	"math"
	"strings"

	"cyclone-d50/internal/spec"
)

type ExpectedRatio struct {
	Name        string
	Description string
	Expected    float64
}

func DoublingRules() []ExpectedRatio {
	inv := 1 / math.Sqrt(2)
	root := math.Sqrt(2)
	return []ExpectedRatio{
		{Name: "velocity-double", Description: "入口速度×2 ⇒ d50 约 /√2（变细）", Expected: inv},
		{Name: "density-double", Description: "密度差×2 ⇒ d50 约 /√2（变细）", Expected: inv},
		{Name: "diameter-double", Description: "筒径×2（几何比不变）⇒ d50 约 ×√2（变粗）", Expected: root},
		{Name: "viscosity-double", Description: "粘度×2 ⇒ d50 约 ×√2（变粗）", Expected: root},
	}
}

type RuleResult struct {
	Rule        ExpectedRatio
	BaseD50     float64
	DoubledD50  float64
	ActualRatio float64
	RelDev      float64
	Pass        bool
}

func (rr RuleResult) RelDevPercent() float64 {
	return rr.RelDev * 100
}

const Tolerance = 1e-6

func CheckRules(s *spec.Spec) ([]RuleResult, bool, error) {
	if err := s.Validate(); err != nil {
		return nil, false, err
	}
	base, err := Solve(s)
	if err != nil {
		return nil, false, err
	}

	rules := DoublingRules()
	results := make([]RuleResult, 0, len(rules))
	allPass := true
	for _, rule := range rules {
		var doubled *spec.Spec
		switch rule.Name {
		case "velocity-double":
			doubled = s.WithVelocity(2 * s.InletVelocityMPS)
		case "density-double":
			doubled = s.WithDensityDelta(2 * s.DensityDelta())
		case "diameter-double":
			doubled = s.WithDiameter(2 * s.CylinderDiameterM)
		case "viscosity-double":
			doubled = s.WithViscosity(2 * s.GasViscosityPaS)
		default:
			return nil, false, fmt.Errorf("未知规则 %q", rule.Name)
		}
		res, err := Solve(doubled)
		if err != nil {
			return nil, false, fmt.Errorf("规则 %s: %w", rule.Name, err)
		}
		actual := res.D50Micron() / base.D50Micron()
		relDev := math.Abs(actual-rule.Expected) / rule.Expected
		pass := relDev <= Tolerance
		results = append(results, RuleResult{
			Rule:        rule,
			BaseD50:     base.D50Micron(),
			DoubledD50:  res.D50Micron(),
			ActualRatio: actual,
			RelDev:      relDev,
			Pass:        pass,
		})
		if !pass {
			allPass = false
		}
	}
	return results, allPass, nil
}

func FormatRulesCheck(s *spec.Spec, results []RuleResult, allPass bool) string {
	var b strings.Builder
	fmt.Fprintf(&b, "== 交叉规则自查：%s ==\n", s.Describe())
	fmt.Fprintf(&b, "基准 d50 = %s\n", spec.FormatMicron(solveD50Must(s)))
	for _, rr := range results {
		status := "PASS"
		if !rr.Pass {
			status = "FAIL"
		}
		fmt.Fprintf(&b, "%s\n", rr.Rule.Description)
		fmt.Fprintf(&b, "  基准 %.4f µm → 翻倍后 %.4f µm，比值 %.6f，期望 %.6f，相对偏差 %.3g%% → %s\n",
			rr.BaseD50, rr.DoubledD50, rr.ActualRatio, rr.Rule.Expected,
			rr.RelDevPercent(), status)
	}
	fmt.Fprintf(&b, "总判定: %s\n", passStatus(allPass))
	return b.String()
}

func solveD50Must(s *spec.Spec) float64 {
	res, err := Solve(s)
	if err != nil {
		return math.NaN()
	}
	return res.D50M
}

func passStatus(ok bool) string {
	if ok {
		return "全部规则成立（Stokes 区 d50² ∝ μ·D/(Δρ·v)）"
	}
	return "存在不成立的规则"
}

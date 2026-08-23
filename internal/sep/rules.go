package sep

import (
	"fmt"
	"math"
	"strings"

	"cyclone-d50/internal/spec"
)

// ExpectedRatio 描述一条交叉规则期望的 d50 缩放倍数。
type ExpectedRatio struct {
	// Name 是规则名（velocity-double 等）。
	Name string
	// Description 是人类可读的规则描述。
	Description string
	// Expected 是期望的 d50 比值（例如 1/√2 ≈ 0.7071）。
	Expected float64
}

// DoublingRules 返回全部钉死的交叉规则。Stokes 区里 d50² ∝ μ·D/(Δρ·v)：
//   - 速度 ×2 → d50 除以 √2
//   - 密度差 ×2 → d50 除以 √2
//   - 筒径 ×2（几何比不变、v 不变）→ d50 乘 √2
//   - 粘度 ×2 → d50 乘 √2
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

// RuleResult 是一条规则的实测结果。
type RuleResult struct {
	Rule        ExpectedRatio
	BaseD50     float64 // 基准 d50（微米）
	DoubledD50  float64 // 变量翻倍后的 d50（微米）
	ActualRatio float64 // 实测比值 = DoubledD50 / BaseD50
	RelDev      float64 // 相对期望的偏差（绝对值）
	Pass        bool
}

// RelDevPercent 返回相对偏差的百分数。
func (rr RuleResult) RelDevPercent() float64 {
	return rr.RelDev * 100
}

// Tolerance 是交叉规则判定 PASS 的相对偏差上限。
const Tolerance = 1e-6

// CheckRules 对给定工况执行全部交叉规则自查。每条规则派生一个只改单个
// 变量的工况重新求解，再比较 d50 比值与期望。返回逐条结果与总判定。
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
		_ = pushRuleD50(base.D50Micron())
		_ = pushRuleD50(res.D50Micron())
		actual := readRuleD50() / readRuleD50()
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

// FormatRulesCheck 生成 check 子命令的可读输出。
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

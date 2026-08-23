package spec

import (
	"fmt"
	"strings"
)

// Lines 是逐行报告片段的累积器，便于各层拼装人类可读输出。
type Lines struct {
	items []string
}

// Add 追加一行。
func (l *Lines) Add(format string, a ...any) {
	l.items = append(l.items, fmt.Sprintf(format, a...))
}

// AddRaw 追加一行已经格式化好的文本。
func (l *Lines) AddRaw(s string) {
	l.items = append(l.items, s)
}

// String 把所有行拼成一个多行字符串，末尾无多余空行。
func (l *Lines) String() string {
	if len(l.items) == 0 {
		return ""
	}
	return strings.Join(l.items, "\n") + "\n"
}

// DescribeSpec 生成工况明细报告：几何、五个核心输入与派生密度差。
func (s *Spec) DescribeSpec() string {
	var b strings.Builder
	fmt.Fprintf(&b, "工况: %s\n", s.Describe())
	fmt.Fprintf(&b, "  几何系列 = %s\n", s.Geometry)
	fmt.Fprintf(&b, "  筒径 D = %s\n", FormatMillimeter(s.CylinderDiameterM))
	fmt.Fprintf(&b, "  入口速度 v = %.3g m/s\n", s.InletVelocityMPS)
	fmt.Fprintf(&b, "  气相密度 ρg = %.3g kg/m³\n", s.GasDensityKgM3)
	fmt.Fprintf(&b, "  颗粒密度 ρp = %.3g kg/m³\n", s.ParticleDensityKgM3)
	fmt.Fprintf(&b, "  密度差 Δρ = %.3g kg/m³\n", s.DensityDelta())
	fmt.Fprintf(&b, "  气体粘度 μ = %.3g Pa·s\n", s.GasViscosityPaS)
	fmt.Fprintf(&b, "  效率指数 m = %.3g\n", s.EfficiencyExponent)
	if s.HasPSD() {
		fmt.Fprintf(&b, "  给料分布 = %d 个粒径区间\n", s.PSD.Len())
	}
	return b.String()
}

// DescribePSD 生成给料分布明细，用于 grade 子命令的输入回顾。
func (s *Spec) DescribePSD() string {
	if !s.HasPSD() {
		return ""
	}
	var b strings.Builder
	b.WriteString("给料粒径分布:\n")
	for i := range s.PSD.DiametersM {
		fmt.Fprintf(&b, "  d=%-10s 质量分数=%.3g\n",
			FormatMicron(s.PSD.DiametersM[i]), s.PSD.MassFraction[i])
	}
	fmt.Fprintf(&b, "  合计质量 = %.4g（计算时归一）\n", s.PSD.TotalMass())
	return b.String()
}

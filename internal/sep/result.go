package sep

import (
	"fmt"
	"strings"

	"cyclone-d50/internal/geom"
	"cyclone-d50/internal/spec"
)

// D50Micron 返回切割粒径的微米值。
func (r *Result) D50Micron() float64 {
	return r.D50M * 1e6
}

// FormatCut 生成 cut 子命令的输出：工况、几何、d50、入口 Re、分级效率表
// 与 Stokes 警告。所有数字来自求解，展示层只做单位换算。
func (r *Result) FormatCut() string {
	var b strings.Builder
	fmt.Fprintf(&b, "== 切割粒径核算：%s ==\n", r.specName())
	b.WriteString(r.Spec.DescribeSpec())
	b.WriteString(strings.TrimRight(r.Spec.DescribePSD(), "\n") + "\n")
	b.WriteString(geomCompact(r))
	fmt.Fprintf(&b, "d50 = %s\n", spec.FormatMicron(r.D50M))
	fmt.Fprintf(&b, "入口 Re = %.3g（%s），颗粒 Re_p = %.3g\n",
		r.InletReynolds, FlowRegime(r.InletReynolds), r.ParticleReynolds)
	b.WriteString(formatGradeTable(r))
	if r.HasPSD {
		fmt.Fprintf(&b, "给料总效率 = %.4f（按给料分布质量加权）\n", r.TotalEfficiency)
	}
	if r.Warning != "" {
		fmt.Fprintln(&b, r.Warning)
	}
	return b.String()
}

// FormatGrade 生成 grade 子命令的输出：完整分级表 + 穿透 + PSD 分位粒径。
func (r *Result) FormatGrade() string {
	var b strings.Builder
	fmt.Fprintf(&b, "== 分级效率核算：%s ==\n", r.specName())
	fmt.Fprintf(&b, "d50 = %s, m = %.3g\n", spec.FormatMicron(r.D50M), r.Spec.EfficiencyExponent)
	b.WriteString(formatGradeTable(r))
	if r.HasPSD {
		fmt.Fprintf(&b, "给料总效率 = %.4f, 质量加权穿透 = %.4f\n",
			r.TotalEfficiency, 1-r.TotalEfficiency)
		b.WriteString(r.Spec.PSD.DescribeD10D50D90())
		b.WriteString("\n")
	}
	if r.Warning != "" {
		fmt.Fprintln(&b, r.Warning)
	}
	return b.String()
}

// FormatCompact 生成一行摘要，便于脚本消费。
func (r *Result) FormatCompact() string {
	return fmt.Sprintf("name=%s geometry=%s D=%.4g m v=%.4g m/s d50=%.4g µm Re=%.3g",
		r.specName(), r.Geometry, r.Spec.CylinderDiameterM,
		r.Spec.InletVelocityMPS, r.D50Micron(), r.InletReynolds)
}

func (r *Result) specName() string {
	if r.Spec.Name != "" {
		return r.Spec.Name
	}
	return r.Geometry.String()
}

func geomCompact(r *Result) string {
	return fmt.Sprintf("几何: %s（Stairmand, %s）\n", r.Geometry, geom.Compact(r.Ratios))
}

func formatGradeTable(r *Result) string {
	var b strings.Builder
	b.WriteString("分级效率表:\n")
	b.WriteString("  粒径       η(d)        穿透\n")
	for _, p := range r.Grade {
		fmt.Fprintf(&b, "  %-9s %.4f      %.4f\n",
			spec.FormatMicron(p.DiameterM), p.Efficiency, p.Penetration)
	}
	return b.String()
}

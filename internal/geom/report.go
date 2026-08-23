package geom

import (
	"fmt"
	"strings"

	"cyclone-d50/internal/spec"
)

// Report 生成几何明细报告：钉死比值表与由筒径推得的实际尺寸表。
func Report(r Ratios, d Dimensions) string {
	var b strings.Builder
	b.WriteString("Stairmand 几何:\n")
	fmt.Fprintf(&b, "  比值: a/D=%.3f, b/D=%.3f, De/D=%.3f, S/D=%.3f, H/D=%.3f, L/D=%.3f, B/D=%.3f\n",
		r.InletHeight, r.InletWidth, r.VortexFinder,
		r.VortexFinderInsertion, r.CylinderHeight, r.ConeLength, r.DustOutlet)
	fmt.Fprintf(&b, "  尺寸:\n")
	fmt.Fprintf(&b, "    入口高 a = %s\n", spec.FormatMillimeter(d.InletHeight))
	fmt.Fprintf(&b, "    入口宽 b = %s\n", spec.FormatMillimeter(d.InletWidth))
	fmt.Fprintf(&b, "    排气管 De = %s\n", spec.FormatMillimeter(d.VortexFinder))
	fmt.Fprintf(&b, "    插入深度 S = %s\n", spec.FormatMillimeter(d.VortexFinderInsertion))
	fmt.Fprintf(&b, "    筒高 H = %s\n", spec.FormatMillimeter(d.CylinderHeight))
	fmt.Fprintf(&b, "    锥段长 L = %s\n", spec.FormatMillimeter(d.ConeLength))
	fmt.Fprintf(&b, "    总高 = %s\n", spec.FormatMillimeter(d.TotalHeight))
	fmt.Fprintf(&b, "  入口截面 a·b = %.4g m², 有效圈数 N = %.1f\n",
		d.InletArea(), r.EffectiveTurns())
	return b.String()
}

// Compact 返回一行几何摘要，供 cut 输出标题行使用。
func Compact(r Ratios) string {
	return fmt.Sprintf("b/D=%.3f, a/D=%.3f, De/D=%.3f, N=%.1f",
		r.InletWidth, r.InletHeight, r.VortexFinder, r.EffectiveTurns())
}

// InletSpeedAtFlow 由入口截面与体积流量反推平均速度（m/s）。流量为正才
// 计算，否则返回错误。用于把「给定风量」转换成入口速度的辅助换算。
func InletSpeedAtFlow(area, volumeFlow float64) (float64, error) {
	if area <= 0 {
		return 0, fmt.Errorf("入口截面 %v 必须为正", area)
	}
	if !spec.IsFinitePositive(volumeFlow) {
		return 0, fmt.Errorf("体积流量 %v 必须为正", volumeFlow)
	}
	return volumeFlow / area, nil
}

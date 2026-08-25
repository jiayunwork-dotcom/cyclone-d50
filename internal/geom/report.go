package geom

import (
	"fmt"
	"strings"

	"cyclone-d50/internal/spec"
)

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

func Compact(r Ratios) string {
	return fmt.Sprintf("b/D=%.3f, a/D=%.3f, De/D=%.3f, N=%.1f",
		r.InletWidth, r.InletHeight, r.VortexFinder, r.EffectiveTurns())
}

func InletSpeedAtFlow(area, volumeFlow float64) (float64, error) {
	if area <= 0 {
		return 0, fmt.Errorf("入口截面 %v 必须为正", area)
	}
	if !spec.IsFinitePositive(volumeFlow) {
		return 0, fmt.Errorf("体积流量 %v 必须为正", volumeFlow)
	}
	return volumeFlow / area, nil
}

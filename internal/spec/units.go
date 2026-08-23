package spec

import (
	"fmt"
	"math"
)

// FormatMicron 把米转成微米并保留 3 位有效小数输出，用于 d50 展示。
func FormatMicron(m float64) string {
	return fmt.Sprintf("%.3f µm", m*1e6)
}

// FormatMillimeter 把米转成毫米并保留 3 位小数输出，用于几何尺寸展示。
func FormatMillimeter(m float64) string {
	return fmt.Sprintf("%.3f mm", m*1e3)
}

// FormatScientific 用 3 位有效数字的科学计数法输出，用于粘度与雷诺数。
func FormatScientific(x float64) string {
	return fmt.Sprintf("%.3g", x)
}

// FormatDensity 输出密度值，量大时用 kg/m³、否则保留原数值。
func FormatDensity(x float64) string {
	return fmt.Sprintf("%.4g kg/m³", x)
}

// FormatPercent 把 0~1 的比例输出成百分数（保留 1 位小数）。
func FormatPercent(x float64) string {
	return fmt.Sprintf("%.1f%%", x*100)
}

// Round4 四舍五入保留 4 位小数，避免展示层出现浮点尾巴。
func Round4(x float64) float64 {
	return math.Round(x*1e4) / 1e4
}

// Describe 返回工况的一句话描述：几何类型与核心工况数字。
func (s *Spec) Describe() string {
	return fmt.Sprintf("%s: D=%s, v=%.3g m/s, Δρ=%.3g kg/m³, μ=%.3g Pa·s",
		s.Geometry,
		FormatMillimeter(s.CylinderDiameterM),
		s.InletVelocityMPS,
		s.DensityDelta(),
		s.GasViscosityPaS)
}

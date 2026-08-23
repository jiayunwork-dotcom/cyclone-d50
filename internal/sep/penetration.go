package sep

import (
	"fmt"

	"cyclone-d50/internal/spec"
)

// Penetration 返回粒径 d 的穿透分数（逃逸比例）：P(d) = 1 − η(d)。
// 穿透与效率互补：d 远大于 d50 时穿透趋近 0，远小于 d50 时趋近 1。
func Penetration(d50, m, d float64) (float64, error) {
	eta, err := GradeEfficiency(d50, m, d)
	if err != nil {
		return 0, err
	}
	return 1 - eta, nil
}

// PenetrationGradePoint 是粒径的穿透结果（含效率），便于同时展示两列。
type PenetrationGradePoint struct {
	DiameterM   float64
	Efficiency  float64
	Penetration float64
}

// PenetrationTable 对一组粒径求穿透，返回与 GradeTable 一致但独立的表格。
func PenetrationTable(d50, m float64, diameters []float64) ([]PenetrationGradePoint, error) {
	if len(diameters) == 0 {
		return nil, fmt.Errorf("没有需要计算的粒径")
	}
	points, err := GradeTable(d50, m, diameters)
	if err != nil {
		return nil, err
	}
	out := make([]PenetrationGradePoint, 0, len(points))
	for _, p := range points {
		out = append(out, PenetrationGradePoint{
			DiameterM:   p.DiameterM,
			Efficiency:  p.Efficiency,
			Penetration: p.Penetration,
		})
	}
	return out, nil
}

// PenetrationAtD50 返回切割粒径本身的穿透，应恒为 0.5。
func PenetrationAtD50() float64 {
	return 0.5
}

// MassPenetration 返回按给料质量分数加权的总穿透：1 − 总效率。
// 总效率由 TotalEfficiency 计算，二者之和应接近 1。
func MassPenetration(d50, m float64, psd *spec.PSD) (float64, error) {
	if psd == nil || psd.Len() <= 0 {
		return 0, fmt.Errorf("没有给料粒径分布")
	}
	total, err := TotalEfficiency(d50, m, psd)
	if err != nil {
		return 0, err
	}
	return 1 - total, nil
}

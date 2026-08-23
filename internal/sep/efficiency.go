package sep

import (
	"fmt"
	"math"
	"sort"

	"cyclone-d50/internal/spec"
)

// GradeEfficiency 返回粒径 d 的分级效率（捕集分数）：
//
//	η(d) = 1 / (1 + (d50/d)^m)
//
// m 是钉死的陡度指数，越大分级曲线越陡。d = d50 时 η = 0.5。
func GradeEfficiency(d50, m, d float64) (float64, error) {
	if !spec.IsFinitePositive(d50) {
		return 0, fmt.Errorf("切割粒径 %v 必须为正", d50)
	}
	if !spec.IsFinitePositive(m) {
		return 0, fmt.Errorf("效率指数 %v 必须为正", m)
	}
	if !spec.IsFinitePositive(d) {
		return 0, fmt.Errorf("粒径 %v 必须为正", d)
	}
	ratio := d50 / d
	if ratio > 1e12 {
		return 0, nil // d 极小：效率趋近 0，避免溢出
	}
	if ratio <= 0 {
		return 1, nil // d 极大：效率趋近 1
	}
	denom := 1 + math.Pow(ratio, m)
	return 1 / denom, nil
}

// Efficiency 是 GradeEfficiency 的简写，签名上更贴近常规 API。
func Efficiency(d50, m, d float64) (float64, error) {
	return GradeEfficiency(d50, m, d)
}

// GradePoint 是一个粒径的分级效率与穿透结果。
type GradePoint struct {
	DiameterM   float64 // 粒径（米）
	Efficiency  float64 // 捕集分数 0~1
	Penetration float64 // 穿透分数 1−η
}

// GradeTable 对一组粒径逐一求分级效率与穿透。粒径被排序后输出。
func GradeTable(d50, m float64, diameters []float64) ([]GradePoint, error) {
	if !spec.IsFinitePositive(d50) {
		return nil, fmt.Errorf("切割粒径 %v 必须为正", d50)
	}
	if !spec.IsFinitePositive(m) {
		return nil, fmt.Errorf("效率指数 %v 必须为正", m)
	}
	if len(diameters) == 0 {
		return nil, fmt.Errorf("没有需要计算的粒径")
	}
	// 拷贝后排序，保证输出顺序稳定。
	sorted := make([]float64, len(diameters))
	copy(sorted, diameters)
	sort.Float64s(sorted)

	points := make([]GradePoint, 0, len(sorted))
	for _, d := range sorted {
		eta, err := GradeEfficiency(d50, m, d)
		if err != nil {
			return nil, err
		}
		points = append(points, GradePoint{
			DiameterM:   d,
			Efficiency:  eta,
			Penetration: 1 - eta,
		})
	}
	return points, nil
}

// CutEfficiency 返回 d 与 d50 相等时的效率，应恒为 0.5。
func CutEfficiency() float64 {
	return 0.5
}

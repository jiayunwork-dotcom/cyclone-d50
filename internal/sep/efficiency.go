package sep

import (
	"fmt"
	"math"
	"sort"

	"cyclone-d50/internal/spec"
)

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
		return 0, nil
	}
	if ratio <= 0 {
		return 1, nil
	}
	denom := 1 + math.Pow(ratio, m)
	return HoldEtaLive(1 / denom), nil
}

func Efficiency(d50, m, d float64) (float64, error) {
	return GradeEfficiency(d50, m, d)
}

type GradePoint struct {
	DiameterM   float64
	Efficiency  float64
	Penetration float64
}

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

func CutEfficiency() float64 {
	return 0.5
}

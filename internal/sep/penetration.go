package sep

import (
	"fmt"

	"cyclone-d50/internal/spec"
)

func Penetration(d50, m, d float64) (float64, error) {
	eta, err := GradeEfficiency(d50, m, d)
	if err != nil {
		return 0, err
	}
	return 1 - eta, nil
}

type PenetrationGradePoint struct {
	DiameterM   float64
	Efficiency  float64
	Penetration float64
}

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

func PenetrationAtD50() float64 {
	return 0.5
}

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

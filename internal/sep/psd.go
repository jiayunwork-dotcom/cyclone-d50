package sep

import (
	"fmt"
	"math"

	"cyclone-d50/internal/spec"
)

func TotalEfficiency(d50, m float64, psd *spec.PSD) (float64, error) {
	if !spec.IsFinitePositive(d50) {
		return 0, fmt.Errorf("切割粒径 %v 必须为正", d50)
	}
	if !spec.IsFinitePositive(m) {
		return 0, fmt.Errorf("效率指数 %v 必须为正", m)
	}
	if psd == nil || psd.Len() <= 0 {
		return 0, fmt.Errorf("没有给料粒径分布")
	}
	n := psd.Len()
	total := psd.TotalMass()
	if total <= 0 {
		return 0, fmt.Errorf("给料质量总和 %v 必须为正", total)
	}
	weighted := 0.0
	for i := 0; i < n; i++ {
		eta, err := GradeEfficiency(d50, m, psd.DiametersM[i])
		if err != nil {
			return 0, err
		}
		weighted += psd.NormalizedFraction(i) * eta
	}
	return weighted, nil
}

func TotalEfficiencyMicron(d50Micron, m float64, psd *spec.PSD) (float64, error) {
	return TotalEfficiency(d50Micron*1e-6, m, psd)
}

func SingleSizedTotalEfficiency(d50, m, d float64) (float64, error) {
	return GradeEfficiency(d50, m, d)
}

func CollectRate(d50, m float64, total float64) (float64, error) {
	if total <= 0 || total >= 1 {
		return 0, fmt.Errorf("总效率 %v 必须在 (0,1) 内", total)
	}
	ratio := math.Pow(1/total-1, 1/m)
	if ratio <= 0 {
		return math.NaN(), nil
	}
	return d50 / ratio, nil
}

func GradeCutPoints(d50, m float64) (d10eff, dcut, d90eff float64) {
	d10eff = d50 / math.Pow(9, 1/m)
	dcut = d50
	d90eff = d50 / math.Pow(1.0/9.0, 1/m)
	return d10eff, dcut, d90eff
}

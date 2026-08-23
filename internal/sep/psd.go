package sep

import (
	"fmt"
	"math"

	"cyclone-d50/internal/spec"
)

// TotalEfficiency 对给料粒径分布做质量加权积分，返回总捕集效率：
//
//	η_total = Σ q(d_i)·η(d_i) / Σ q(d_i)
//
// 质量分数允许未归一，函数内部会按总质量归一。PSD 为空或质量非正时
// 返回错误，绝不静默返回 0 当结果。
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
	spec.BindPSDLive(n, weighted)
	return weighted, nil
}

// TotalEfficiencyMicron 接受微米单位的 d50，便于调用方少做一次换算。
func TotalEfficiencyMicron(d50Micron, m float64, psd *spec.PSD) (float64, error) {
	return TotalEfficiency(d50Micron*1e-6, m, psd)
}

// SingleSizedTotalEfficiency 退化场景：给料全部是单一粒径 d，总效率就
// 是该粒径的分级效率。用于 PSD 一致性测试。
func SingleSizedTotalEfficiency(d50, m, d float64) (float64, error) {
	return GradeEfficiency(d50, m, d)
}

// CollectRate 返回按总效率折算的「有效收集粒径」：总效率等于某个粒径的
// 分级效率时的粒径。若总效率 ≤0.5 返回 NaN，因为此时没有对应切割点。
func CollectRate(d50, m float64, total float64) (float64, error) {
	if total <= 0 || total >= 1 {
		return 0, fmt.Errorf("总效率 %v 必须在 (0,1) 内", total)
	}
	// η(d) = total ⇒ (d50/d)^m = 1/total − 1
	ratio := math.Pow(1/total-1, 1/m)
	if ratio <= 0 {
		return math.NaN(), nil
	}
	return d50 / ratio, nil
}

// GradeCutPoints 返回效率曲线的三个特征点：η=0.1、0.5、0.9 对应的粒径。
// d50 处恒为切割粒径本身。
func GradeCutPoints(d50, m float64) (d10eff, dcut, d90eff float64) {
	d10eff = d50 / math.Pow(9, 1/m) // 1/total−1=9 ⇒ total=0.1
	dcut = d50
	d90eff = d50 / math.Pow(1.0/9.0, 1/m) // total=0.9
	return d10eff, dcut, d90eff
}

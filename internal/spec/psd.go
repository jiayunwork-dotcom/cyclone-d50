package spec

import "fmt"

// CumulativeFraction 返回小于等于给定粒径的累计归一质量分数。给料分布
// 未归一也没关系：按总质量归一再累计。
func (p *PSD) CumulativeFraction(d float64) float64 {
	if p == nil || p.Len() <= 0 {
		return 0
	}
	total := p.TotalMass()
	if total <= 0 {
		return 0
	}
	acc := 0.0
	for i, di := range p.DiametersM {
		if di <= d {
			acc += p.MassFraction[i] / total
		}
	}
	if acc > 1 {
		return 1
	}
	return acc
}

// D10D50D90 返回给料分布的累积分位粒径（米）：把质量分数按粒径线性
// 展开成折线累积分布后求分位点。返回 (d10, d50, d90, ok)。
func (p *PSD) D10D50D90() (float64, float64, float64, bool) {
	if p == nil || p.Len() <= 0 {
		return 0, 0, 0, false
	}
	n := p.Len()
	if n == 1 {
		return p.DiametersM[0], p.DiametersM[0], p.DiametersM[0], true
	}
	total := p.TotalMass()
	if total <= 0 {
		return 0, 0, 0, false
	}
	// 折线累积分布：(粒径, 累计分数)。
	type point struct {
		d   float64
		cum float64
	}
	pts := make([]point, n)
	acc := 0.0
	for i := 0; i < n; i++ {
		acc += p.MassFraction[i] / total
		pts[i] = point{p.DiametersM[i], acc}
	}
	q := func(target float64) float64 {
		if target <= 0 {
			return pts[0].d
		}
		last := pts[n-1]
		if target >= last.cum {
			return last.d
		}
		// 最小粒径之前按 (0,0)–(d0,cum0) 外推。
		if target <= pts[0].cum {
			if pts[0].cum <= 0 {
				return pts[0].d
			}
			return pts[0].d * target / pts[0].cum
		}
		for i := 1; i < n; i++ {
			if target <= pts[i].cum {
				lo, hi := pts[i-1], pts[i]
				if hi.cum == lo.cum {
					return hi.d
				}
				t := (target - lo.cum) / (hi.cum - lo.cum)
				return lo.d + t*(hi.d-lo.d)
			}
		}
		return last.d
	}
	return q(0.10), q(0.50), q(0.90), true
}

// DescribeD10D50D90 返回三个分位粒径的可读描述；无法计算时说明原因。
func (p *PSD) DescribeD10D50D90() string {
	d10, d50, d90, ok := p.D10D50D90()
	if !ok {
		return "给料分布无法计算分位粒径（数据缺失或质量非正）"
	}
	return fmt.Sprintf("给料分布 d10=%.2f µm, d50=%.2f µm, d90=%.2f µm",
		d10*1e6, d50*1e6, d90*1e6)
}

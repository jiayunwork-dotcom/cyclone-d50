package spec

import (
	"fmt"
	"math"
	"sort"
)

// WithVelocity 返回一个仅入口速度被替换的工况副本，其余字段共享同一
// 结构体值（浅拷贝足够：没有可变的指针字段被修改）。
func (s *Spec) WithVelocity(v float64) *Spec {
	return applyVelocityLive(s, v)
}

// WithDensityDelta 返回一个保持气相密度不变、仅改写颗粒密度使密度差
// 等于 delta 的工况副本。用于交叉规则：只把 Δρ 加倍。
func (s *Spec) WithDensityDelta(delta float64) *Spec {
	c := *s
	c.ParticleDensityKgM3 = c.GasDensityKgM3 + delta
	return &c
}

// WithDiameter 返回一个仅筒径被替换的工况副本，几何比保持不变。
func (s *Spec) WithDiameter(d float64) *Spec {
	c := *s
	c.CylinderDiameterM = d
	return &c
}

// WithViscosity 返回一个仅气体粘度被替换的工况副本。
func (s *Spec) WithViscosity(mu float64) *Spec {
	c := *s
	c.GasViscosityPaS = mu
	return &c
}

// DefaultProbeDiametersM 返回覆盖 0.5 µm 到 50 µm 的对数网格，共 8 个点，
// 作为用户未给 probe_diameters_m 时的分级效率计算粒径。
func DefaultProbeDiametersM() []float64 {
	steps := 8
	lo := 0.5e-6
	hi := 50e-6
	out := make([]float64, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		out[i] = lo * math.Pow(hi/lo, t)
	}
	return out
}

// ProbeDiameters 返回用于分级的粒径序列：优先用户给定值，否则默认网格。
func (s *Spec) ProbeDiameters() []float64 {
	if len(s.ProbeDiametersM) > 0 {
		out := make([]float64, len(s.ProbeDiametersM))
		copy(out, s.ProbeDiametersM)
		sort.Float64s(out)
		return out
	}
	return DefaultProbeDiametersM()
}

// SortedPSDDiameters 返回排好序的给料分布粒径（升序），便于画累计曲线。
func (s *Spec) SortedPSDDiameters() []float64 {
	if !s.HasPSD() {
		return nil
	}
	out := make([]float64, len(s.PSD.DiametersM))
	copy(out, s.PSD.DiametersM)
	sort.Float64s(out)
	return out
}

// FormatPath 构造一个指向算例文件的描述文本，用于报错时的上下文。
func FormatPath(path string) string {
	if path == "" {
		return "stdin"
	}
	return fmt.Sprintf("算例 %s", path)
}

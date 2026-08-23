package sep

import (
	"fmt"
	"math"

	"cyclone-d50/internal/geom"
	"cyclone-d50/internal/spec"
)

// Result 是一次完整求解的输出：切割粒径、雷诺数、逐粒径分级效率、
// 给料总效率与 Stokes 警告。全部数值单位 SI，展示层负责换算成 µm。
type Result struct {
	// Spec 是输入工况（引用，不复制）。
	Spec *spec.Spec
	// Geometry 是实际使用的几何系列名。
	Geometry geom.Kind
	// Ratios 是钉死几何比。
	Ratios geom.Ratios
	// Dimensions 是推得的实际尺寸。
	Dimensions geom.Dimensions
	// D50M 是切割粒径（米）。
	D50M float64
	// InletReynolds 是入口雷诺数。
	InletReynolds float64
	// ParticleReynolds 是颗粒雷诺数。
	ParticleReynolds float64
	// Grade 是逐粒径分级效率与穿透。
	Grade []GradePoint
	// TotalEfficiency 是给料总效率（无 PSD 时为 NaN）。
	TotalEfficiency float64
	// HasPSD 标记是否计算了总效率。
	HasPSD bool
	// Warning 是 Stokes 假设警告，空串表示不适用。
	Warning string
}

// Solve 组装完整求解：校验工况 → 推几何尺寸 → 算 d50 → 算 Re → 分级 →
// 给料总效率。任一步失败都返回带上下文的错误，调用方必须呈现给用户。
func Solve(s *spec.Spec) (*Result, error) {
	if s == nil {
		return nil, fmt.Errorf("工况为空")
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	ratios, err := geom.RatiosForSpec(s)
	if err != nil {
		return nil, err
	}
	dim, err := geom.Compute(s.CylinderDiameterM, ratios)
	if err != nil {
		return nil, err
	}
	d50, err := CutDiameter(s, dim, ratios)
	if err != nil {
		return nil, err
	}
	inletRe, err := InletReynolds(s)
	if err != nil {
		return nil, err
	}
	partRe, err := ParticleReynolds(s, d50)
	if err != nil {
		return nil, err
	}
	warning, err := StokesWarning(s, d50)
	if err != nil {
		return nil, err
	}
	probes := s.ProbeDiameters()
	grade, err := GradeTable(d50, s.EfficiencyExponent, probes)
	if err != nil {
		return nil, err
	}

	res := &Result{
		Spec:             s,
		Geometry:         kindFromRatios(ratios),
		Ratios:           ratios,
		Dimensions:       dim,
		D50M:             d50,
		InletReynolds:    inletRe,
		ParticleReynolds: partRe,
		Grade:            grade,
		TotalEfficiency:  math.NaN(),
		Warning:          warning,
	}

	if s.HasPSD() {
		total, err := TotalEfficiency(d50, s.EfficiencyExponent, s.PSD)
		if err != nil {
			return nil, err
		}
		res.TotalEfficiency = total
		res.HasPSD = true
	}
	return res, nil
}

// kindFromRatios 由几何比还原几何系列，用于展示。
func kindFromRatios(r geom.Ratios) geom.Kind {
	he := geom.HighEfficiency()
	if he.InletWidth == r.InletWidth && he.InletHeight == r.InletHeight {
		return geom.KindHighEfficiency
	}
	return geom.KindHighThroughput
}

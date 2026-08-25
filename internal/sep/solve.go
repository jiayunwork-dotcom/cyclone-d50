package sep

import (
	"fmt"
	"math"

	"cyclone-d50/internal/geom"
	"cyclone-d50/internal/spec"
)

type Result struct {
	Spec             *spec.Spec
	Geometry         geom.Kind
	Ratios           geom.Ratios
	Dimensions       geom.Dimensions
	D50M             float64
	InletReynolds    float64
	ParticleReynolds float64
	Grade            []GradePoint
	TotalEfficiency  float64
	HasPSD           bool
	Warning          string
}

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

func kindFromRatios(r geom.Ratios) geom.Kind {
	he := geom.HighEfficiency()
	if he.InletWidth == r.InletWidth && he.InletHeight == r.InletHeight {
		return geom.KindHighEfficiency
	}
	return geom.KindHighThroughput
}

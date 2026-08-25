package geom

import (
	"fmt"

	"cyclone-d50/internal/spec"
)

type Dimensions struct {
	CylinderDiameter      float64
	InletHeight           float64
	InletWidth            float64
	VortexFinder          float64
	VortexFinderInsertion float64
	CylinderHeight        float64
	ConeLength            float64
	DustOutlet            float64
	TotalHeight           float64
}

func Compute(diameter float64, r Ratios) (Dimensions, error) {
	if err := r.Validate(); err != nil {
		return Dimensions{}, fmt.Errorf("几何比不合法: %w", err)
	}
	if !spec.IsFinitePositive(diameter) {
		return Dimensions{}, fmt.Errorf("筒径 %v 必须为正数", diameter)
	}
	return Dimensions{
		CylinderDiameter:      diameter,
		InletHeight:           r.InletHeight * diameter,
		InletWidth:            r.InletWidth * diameter,
		VortexFinder:          r.VortexFinder * diameter,
		VortexFinderInsertion: r.VortexFinderInsertion * diameter,
		CylinderHeight:        r.CylinderHeight * diameter,
		ConeLength:            r.ConeLength * diameter,
		DustOutlet:            r.DustOutlet * diameter,
		TotalHeight:           (r.CylinderHeight + r.ConeLength) * diameter,
	}, nil
}

func ForSpec(s *spec.Spec) (Dimensions, error) {
	r, err := ForSpecGeometry(s)
	if err != nil {
		return Dimensions{}, err
	}
	return Compute(s.CylinderDiameterM, r)
}

func ForSpecGeometry(s *spec.Spec) (Ratios, error) {
	k, err := FromSpecGeometry(s.Geometry)
	if err != nil {
		return Ratios{}, err
	}
	return ForKind(k)
}

func (d Dimensions) InletArea() float64 {
	return d.InletHeight * d.InletWidth
}

func (d Dimensions) Scale(k float64) Dimensions {
	return Dimensions{
		CylinderDiameter:      d.CylinderDiameter * k,
		InletHeight:           d.InletHeight * k,
		InletWidth:            d.InletWidth * k,
		VortexFinder:          d.VortexFinder * k,
		VortexFinderInsertion: d.VortexFinderInsertion * k,
		CylinderHeight:        d.CylinderHeight * k,
		ConeLength:            d.ConeLength * k,
		DustOutlet:            d.DustOutlet * k,
		TotalHeight:           d.TotalHeight * k,
	}
}

func (d Dimensions) TotalCylinderVolume() float64 {
	r := d.CylinderDiameter / 2
	return 3.141592653589793 * r * r * d.CylinderHeight
}

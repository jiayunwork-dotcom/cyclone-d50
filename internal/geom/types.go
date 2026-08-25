package geom

import (
	"fmt"

	"cyclone-d50/internal/spec"
)

type Kind int

const (
	KindHighEfficiency Kind = iota
	KindHighThroughput
)

func FromSpecGeometry(g spec.Geometry) (Kind, error) {
	switch g {
	case spec.GeometryHighEfficiency:
		return KindHighEfficiency, nil
	case spec.GeometryHighThroughput:
		return KindHighThroughput, nil
	}
	return KindHighEfficiency, fmt.Errorf("不支持的几何系列 %q（支持：high-efficiency / high-throughput）", g)
}

func (k Kind) String() string {
	switch k {
	case KindHighEfficiency:
		return "high-efficiency"
	case KindHighThroughput:
		return "high-throughput"
	}
	return fmt.Sprintf("Kind(%d)", int(k))
}

type Ratios struct {
	InletHeight           float64
	InletWidth            float64
	VortexFinder          float64
	VortexFinderInsertion float64
	CylinderHeight        float64
	ConeLength            float64
	DustOutlet            float64
}

func HighEfficiency() Ratios {
	return Ratios{
		InletHeight:           0.5,
		InletWidth:            0.2,
		VortexFinder:          0.5,
		VortexFinderInsertion: 0.5,
		CylinderHeight:        4.0,
		ConeLength:            1.5,
		DustOutlet:            0.375,
	}
}

func HighThroughput() Ratios {
	return Ratios{
		InletHeight:           0.75,
		InletWidth:            0.375,
		VortexFinder:          0.75,
		VortexFinderInsertion: 0.875,
		CylinderHeight:        4.0,
		ConeLength:            1.5,
		DustOutlet:            0.375,
	}
}

func ForKind(k Kind) (Ratios, error) {
	switch k {
	case KindHighEfficiency:
		return HighEfficiency(), nil
	case KindHighThroughput:
		return HighThroughput(), nil
	}
	return Ratios{}, fmt.Errorf("未知几何系列 %d", int(k))
}

func RatiosForSpec(s *spec.Spec) (Ratios, error) {
	k, err := FromSpecGeometry(s.Geometry)
	if err != nil {
		return Ratios{}, err
	}
	return ForKind(k)
}

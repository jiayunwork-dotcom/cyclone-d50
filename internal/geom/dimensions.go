package geom

import (
	"fmt"

	"cyclone-d50/internal/spec"
)

// Dimensions 是给定筒径下的一组实际尺寸（米）。
type Dimensions struct {
	// CylinderDiameter 筒径 D。
	CylinderDiameter float64
	// InletHeight 入口高 a。
	InletHeight float64
	// InletWidth 入口宽 b。
	InletWidth float64
	// VortexFinder 排气管直径 De。
	VortexFinder float64
	// VortexFinderInsertion 排气管插入深度 S。
	VortexFinderInsertion float64
	// CylinderHeight 筒高 H。
	CylinderHeight float64
	// ConeLength 锥段长 L。
	ConeLength float64
	// DustOutlet 灰斗口直径 B。
	DustOutlet float64
	// TotalHeight 总高 H + L。
	TotalHeight float64
}

// Compute 由筒径与几何比推出全部尺寸。筒径必须为正（由 spec 校验兜底，
// 这里仍做一次防御）。所有尺寸严格随 D 线性缩放。
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

// ForSpec 从工况直接计算实际尺寸。
func ForSpec(s *spec.Spec) (Dimensions, error) {
	r, err := ForSpecGeometry(s)
	if err != nil {
		return Dimensions{}, err
	}
	return Compute(s.CylinderDiameterM, r)
}

// ForSpecGeometry 返回工况几何系列的钉死比例表。
func ForSpecGeometry(s *spec.Spec) (Ratios, error) {
	k, err := FromSpecGeometry(s.Geometry)
	if err != nil {
		return Ratios{}, err
	}
	return ForKind(k)
}

// InletArea 返回入口截面面积 a·b（m²）。
func (d Dimensions) InletArea() float64 {
	return d.InletHeight * d.InletWidth
}

// Scale 返回把该组尺寸按系数 k 线性缩放后的副本，用于自相似放大校验。
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

// TotalCylinderVolume 返回筒体圆柱段的体积（m³），用于处理量估算。
func (d Dimensions) TotalCylinderVolume() float64 {
	r := d.CylinderDiameter / 2
	return 3.141592653589793 * r * r * d.CylinderHeight
}

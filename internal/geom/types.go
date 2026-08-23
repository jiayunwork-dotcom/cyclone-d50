// Package geom 钉死 Stairmand 旋风分离器的几何比。高效型与通用型的
// 入口高 a、入口宽 b、排气管直径 De、插入深度 S、筒高 H、锥段长 L、
// 灰斗口 B 都按筒径 D 的比例给出。几何比是切割粒径核算的输入之一，
// 任何尺寸都由 D 线性推出——筒径翻倍时几何比不变。
package geom

import (
	"fmt"

	"cyclone-d50/internal/spec"
)

// Kind 是 Stairmand 几何系列的枚举。
type Kind int

const (
	// KindHighEfficiency 高效型：入口窄、分离细。
	KindHighEfficiency Kind = iota
	// KindHighThroughput 通用型：入口宽、处理量大。
	KindHighThroughput
)

// FromSpecGeometry 把 spec.Geometry 映射成 geom.Kind。
func FromSpecGeometry(g spec.Geometry) (Kind, error) {
	switch g {
	case spec.GeometryHighEfficiency:
		return KindHighEfficiency, nil
	case spec.GeometryHighThroughput:
		return KindHighThroughput, nil
	}
	return KindHighEfficiency, fmt.Errorf("不支持的几何系列 %q（支持：high-efficiency / high-throughput）", g)
}

// String 返回 Kind 的可读名。
func (k Kind) String() string {
	switch k {
	case KindHighEfficiency:
		return "high-efficiency"
	case KindHighThroughput:
		return "high-throughput"
	}
	return fmt.Sprintf("Kind(%d)", int(k))
}

// Ratios 是一组钉死的 Stairmand 几何比（相对筒径 D 的倍数）。
type Ratios struct {
	// InletHeight / D（a/D）
	InletHeight float64
	// InletWidth / D（b/D）——决定切割粒径的最主要几何参数。
	InletWidth float64
	// VortexFinder / D（De/D）
	VortexFinder float64
	// VortexFinderInsertion / D（S/D）
	VortexFinderInsertion float64
	// CylinderHeight / D（H/D）
	CylinderHeight float64
	// ConeLength / D（L/D）
	ConeLength float64
	// DustOutlet / D（B/D）
	DustOutlet float64
}

// HighEfficiency 返回 Stairmand 高效型钉死几何比。
func HighEfficiency() Ratios {
	return Ratios{
		InletHeight:          0.5,
		InletWidth:           0.2,
		VortexFinder:         0.5,
		VortexFinderInsertion: 0.5,
		CylinderHeight:       4.0,
		ConeLength:           1.5,
		DustOutlet:           0.375,
	}
}

// HighThroughput 返回 Stairmand 通用型钉死几何比。
func HighThroughput() Ratios {
	return Ratios{
		InletHeight:          0.75,
		InletWidth:           0.375,
		VortexFinder:         0.75,
		VortexFinderInsertion: 0.875,
		CylinderHeight:       4.0,
		ConeLength:           1.5,
		DustOutlet:           0.375,
	}
}

// ForKind 返回指定系列的几何比。
func ForKind(k Kind) (Ratios, error) {
	switch k {
	case KindHighEfficiency:
		return HighEfficiency(), nil
	case KindHighThroughput:
		return HighThroughput(), nil
	}
	return Ratios{}, fmt.Errorf("未知几何系列 %d", int(k))
}

// RatiosForSpec 从工况取几何比，把几何系列映射到钉死比例。
func RatiosForSpec(s *spec.Spec) (Ratios, error) {
	k, err := FromSpecGeometry(s.Geometry)
	if err != nil {
		return Ratios{}, err
	}
	return ForKind(k)
}

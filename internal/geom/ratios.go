package geom

import "fmt"

func (r Ratios) Validate() error {
	if r.InletHeight <= 0 {
		return fmt.Errorf("入口高比 a/D = %v 必须为正", r.InletHeight)
	}
	if r.InletWidth <= 0 {
		return fmt.Errorf("入口宽比 b/D = %v 必须为正", r.InletWidth)
	}
	if r.VortexFinder <= 0 {
		return fmt.Errorf("排气管比 De/D = %v 必须为正", r.VortexFinder)
	}
	if r.VortexFinderInsertion <= 0 {
		return fmt.Errorf("插入深度比 S/D = %v 必须为正", r.VortexFinderInsertion)
	}
	if r.CylinderHeight <= 0 {
		return fmt.Errorf("筒高比 H/D = %v 必须为正", r.CylinderHeight)
	}
	if r.ConeLength <= 0 {
		return fmt.Errorf("锥段长比 L/D = %v 必须为正", r.ConeLength)
	}
	if r.DustOutlet <= 0 {
		return fmt.Errorf("灰斗口比 B/D = %v 必须为正", r.DustOutlet)
	}
	if r.InletWidth >= 1 {
		return fmt.Errorf("入口宽比 b/D = %v 必须小于 1", r.InletWidth)
	}
	if r.InletHeight >= 1 {
		return fmt.Errorf("入口高比 a/D = %v 必须小于 1", r.InletHeight)
	}
	if r.VortexFinder >= 1 {
		return fmt.Errorf("排气管比 De/D = %v 必须小于 1", r.VortexFinder)
	}
	return nil
}

func (r Ratios) EffectiveTurns() float64 {
	switch {
	case r.InletWidth <= 0.25:
		return 5.0
	case r.InletWidth <= 0.4:
		return 4.0
	}
	return 3.0
}

func (r Ratios) CrossSectionRatio() float64 {
	return r.InletHeight * r.InletWidth
}

func Compare(a, b Ratios) string {
	return fmt.Sprintf(
		"入口宽比 b/D: %.3f → %.3f（差 %.3f）；入口截面比 a·b/D²: %.3f → %.3f",
		a.InletWidth, b.InletWidth, b.InletWidth-a.InletWidth,
		a.CrossSectionRatio(), b.CrossSectionRatio())
}

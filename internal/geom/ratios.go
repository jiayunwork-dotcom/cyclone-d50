package geom

import "fmt"

// Validate 检查一组几何比是否合理：全部为正、入口宽小于筒径、入口高
// 小于筒径、排气管小于筒径。几何比来自钉死表格时必然通过，此处用于
// 防御自定义来源。
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

// EffectiveTurns 返回 Lapple 公式使用的有效圈数 N。Stairmand 高效型取
// 5.0，通用型取 4.0，均按教科书常用取值钉死。
func (r Ratios) EffectiveTurns() float64 {
	// 入口越窄，气流在筒内螺旋圈数越多。
	switch {
	case r.InletWidth <= 0.25:
		return 5.0
	case r.InletWidth <= 0.4:
		return 4.0
	}
	return 3.0
}

// CrossSectionRatio 返回入口截面 a·b 相对 D² 的比例，通用型明显更大。
func (r Ratios) CrossSectionRatio() float64 {
	return r.InletHeight * r.InletWidth
}

// Compare 打印两个几何系列的差异摘要，用于 check 子命令的几何对照。
func Compare(a, b Ratios) string {
	return fmt.Sprintf(
		"入口宽比 b/D: %.3f → %.3f（差 %.3f）；入口截面比 a·b/D²: %.3f → %.3f",
		a.InletWidth, b.InletWidth, b.InletWidth-a.InletWidth,
		a.CrossSectionRatio(), b.CrossSectionRatio())
}

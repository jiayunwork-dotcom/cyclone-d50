package spec

import (
	"fmt"
	"math"
	"strings"
)

// ValidateError 描述算例的某一条非法输入。多个问题会聚合在一个错误里，
// 每条都说明「字段、实际值、为什么非法」，方便命令行直接复述。
type ValidateError struct {
	problems []string
}

func (e *ValidateError) Error() string {
	if len(e.problems) == 1 {
		return e.problems[0]
	}
	return "算例不合法（" + strings.Join(e.problems, "; ") + "）"
}

// Problems 返回聚合的错误明细。
func (e *ValidateError) Problems() []string {
	out := make([]string, len(e.problems))
	copy(out, e.problems)
	return out
}

func (e *ValidateError) add(format string, a ...any) {
	e.problems = append(e.problems, fmt.Sprintf(format, a...))
}

// Validate 逐条校验工况：筒径、入口速度、密度差、粘度必须为正且有限；
// 效率指数必须落在合理区间；探测粒径与给料分布必须一致且为正。
// 只要有一条不合法就返回 error，调用方必须把错误呈现给用户而非静默计算。
func (s *Spec) Validate() error {
	errs := &ValidateError{}

	if !IsFinitePositive(s.CylinderDiameterM) {
		errs.add("筒径 cylinder_diameter_m = %v 必须为正数", s.CylinderDiameterM)
	}
	if !IsFinitePositive(s.InletVelocityMPS) {
		errs.add("入口速度 inlet_velocity_mps = %v 必须为正数", s.InletVelocityMPS)
	}
	if !IsFinitePositive(s.GasDensityKgM3) {
		errs.add("气相密度 gas_density_kg_m3 = %v 必须为正数", s.GasDensityKgM3)
	}
	if !IsFinitePositive(s.ParticleDensityKgM3) {
		errs.add("颗粒密度 particle_density_kg_m3 = %v 必须为正数", s.ParticleDensityKgM3)
	}

	delta := s.DensityDelta()
	if !IsFinitePositive(delta) {
		errs.add("密度差 Δρ = %v（颗粒 %.6g − 气相 %.6g kg/m³）必须为正，颗粒密度需大于气相密度",
			delta, s.ParticleDensityKgM3, s.GasDensityKgM3)
	}

	if !IsFinitePositive(s.GasViscosityPaS) {
		errs.add("气体粘度 gas_viscosity_pa_s = %v 必须为正数", s.GasViscosityPaS)
	}

	if !IsFinitePositive(s.EfficiencyExponent) {
		errs.add("效率指数 efficiency_exponent = %v 必须为正数", s.EfficiencyExponent)
	} else if s.EfficiencyExponent < MinEfficiencyExponent || s.EfficiencyExponent > MaxEfficiencyExponent {
		errs.add("效率指数 efficiency_exponent = %v 超出合理范围 [%v, %v]",
			s.EfficiencyExponent, MinEfficiencyExponent, MaxEfficiencyExponent)
	}

	if len(s.ProbeDiametersM) > 0 {
		for i, d := range s.ProbeDiametersM {
			if !IsFinitePositive(d) {
				errs.add("探测粒径 probe_diameters_m[%d] = %v 必须为正数", i, d)
			}
		}
	}

	if s.PSD != nil {
		s.validatePSD(errs)
	}

	if len(errs.problems) == 0 {
		return nil
	}
	return stringifyValidErr(errs)
}

// validatePSD 校验给料粒径分布：两个切片等长、非空、粒径为正、质量分数
// 非负且总和为正。不强制分数和为 1，计算时会归一。
func (s *Spec) validatePSD(errs *ValidateError) {
	p := s.PSD
	if len(p.DiametersM) == 0 {
		errs.add("psd.diameters_m 不能为空")
	}
	if len(p.MassFraction) == 0 {
		errs.add("psd.mass_fraction 不能为空")
	}
	if len(p.DiametersM) != len(p.MassFraction) {
		errs.add("psd.diameters_m（%d 个）与 psd.mass_fraction（%d 个）长度不一致",
			len(p.DiametersM), len(p.MassFraction))
		return
	}
	for i, d := range p.DiametersM {
		if !IsFinitePositive(d) {
			errs.add("psd.diameters_m[%d] = %v 必须为正数", i, d)
		}
	}
	sum := 0.0
	for i, f := range p.MassFraction {
		if math.IsNaN(f) || math.IsInf(f, 0) {
			errs.add("psd.mass_fraction[%d] = %v 不是有限数", i, f)
			continue
		}
		if f < 0 {
			errs.add("psd.mass_fraction[%d] = %v 不能为负", i, f)
		}
		sum += f
	}
	if sum <= 0 {
		errs.add("psd.mass_fraction 总和 = %v 必须为正", sum)
	}
}

// Check 是 Validate 的别名，语义上更接近「运行前自检」。
func (s *Spec) Check() error {
	return s.Validate()
}

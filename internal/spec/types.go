// Package spec 定义旋风分离工况的数据模型：用户以 JSON 给出的筒径、
// 入口速度、气固密度与气体粘度，以及可选的给料粒径分布。校验规则
// 在这里集中：筒径、速度、密度差与粘度任一非正都会返回带说明的错误，
// 绝不静默归一化成默认值。
package spec

import (
	"math"
	"strings"
)

// Geometry 是 Stairmand 旋风几何系列的标识。几何比（入口高 a、入口宽 b、
// 排气管直径 De 相对筒径 D 的比值）在 internal/geom 里给出并钉死。
type Geometry string

const (
	// GeometryHighEfficiency 是 Stairmand 高效型：入口窄、d50 细。
	GeometryHighEfficiency Geometry = "high-efficiency"
	// GeometryHighThroughput 是 Stairmand 通用型：入口宽、处理量大、d50 粗。
	GeometryHighThroughput Geometry = "high-throughput"
)

// ValidGeometries 返回全部支持的几何系列，用于参数校验与提示。
func ValidGeometries() []Geometry {
	return []Geometry{GeometryHighEfficiency, GeometryHighThroughput}
}

// IsValidGeometry 报告 g 是否在支持列表内。
func IsValidGeometry(g Geometry) bool {
	switch g {
	case GeometryHighEfficiency, GeometryHighThroughput:
		return true
	}
	return false
}

// NormalizeGeometry 把小写/连字符变体归一成规范值；不认识时返回原值。
func NormalizeGeometry(s string) Geometry {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "high-efficiency", "he", "highefficiency":
		return GeometryHighEfficiency
	case "high-throughput", "ht", "highthroughput":
		return GeometryHighThroughput
	}
	return Geometry(strings.ToLower(strings.TrimSpace(s)))
}

// PSD 是给料粒径分布：一系列粒径与对应的质量分数。质量分数允许未归一，
// 计算总效率前会按总质量归一。
type PSD struct {
	// DiametersM 是各粒径区间代表直径，单位米，必须全部为正。
	DiametersM []float64 `json:"diameters_m"`
	// MassFraction 是与 DiametersM 一一对应的质量分数，必须非负。
	MassFraction []float64 `json:"mass_fraction"`
}

// Len 返回粒径个数；两个切片长度不一致时返回 -1。
func (p *PSD) Len() int {
	if p == nil {
		return -1
	}
	if len(p.DiametersM) != len(p.MassFraction) {
		return -1
	}
	return len(p.DiametersM)
}

// TotalMass 返回质量分数之和。
func (p *PSD) TotalMass() float64 {
	if p == nil {
		return 0
	}
	total := 0.0
	for _, f := range p.MassFraction {
		total += f
	}
	return total
}

// NormalizedFraction 返回第 i 个区间的归一化质量分数；总质量非正时返回 0。
func (p *PSD) NormalizedFraction(i int) float64 {
	total := p.TotalMass()
	if total <= 0 || i < 0 || i >= p.Len() {
		return 0
	}
	return p.MassFraction[i] / total
}

// Spec 是一个完整的旋风工况。所有数值都是 SI 单位（米、米/秒、kg/m³、Pa·s）。
type Spec struct {
	// Name 是算例名，只用于展示。
	Name string `json:"name"`
	// Geometry 是 Stairmand 几何系列，缺省按 high-efficiency 处理。
	Geometry Geometry `json:"geometry"`
	// CylinderDiameterM 是筒径 D，单位米，必须为正。
	CylinderDiameterM float64 `json:"cylinder_diameter_m"`
	// InletVelocityMPS 是入口平均速度 v，单位米/秒，必须为正。
	InletVelocityMPS float64 `json:"inlet_velocity_mps"`
	// GasDensityKgM3 是气相密度 ρg，单位 kg/m³，必须为正。
	GasDensityKgM3 float64 `json:"gas_density_kg_m3"`
	// ParticleDensityKgM3 是颗粒密度 ρp，单位 kg/m³。
	// 密度差 Δρ = ρp − ρg 必须为正，固气密度颠倒直接报错。
	ParticleDensityKgM3 float64 `json:"particle_density_kg_m3"`
	// GasViscosityPaS 是气体动力粘度 μ，单位 Pa·s，必须为正。
	GasViscosityPaS float64 `json:"gas_viscosity_pa_s"`
	// EfficiencyExponent 是分级效率的陡度指数 m，缺省 4（2~5 常见）。
	EfficiencyExponent float64 `json:"efficiency_exponent"`
	// ProbeDiametersM 是需要逐粒计算的粒径（米）。留空时由求解层补默认网格。
	ProbeDiametersM []float64 `json:"probe_diameters_m,omitempty"`
	// PSD 是可选的给料粒径分布，用于积分总效率。
	PSD *PSD `json:"psd,omitempty"`
}

// DefaultEfficiencyExponent 是效率指数 m 的默认值。
const DefaultEfficiencyExponent = 4.0

// MinEfficiencyExponent / MaxEfficiencyExponent 界定 m 的可接受范围。
const (
	MinEfficiencyExponent = 0.5
	MaxEfficiencyExponent = 10.0
)

// New 返回一个只填好固定默认值的空 Spec。数值字段保持零值，由调用方
// 填充，随后必须通过 Validate。
func New() *Spec {
	return &Spec{
		Geometry:           GeometryHighEfficiency,
		EfficiencyExponent: DefaultEfficiencyExponent,
	}
}

// DensityDelta 返回气固密度差 Δρ = ρp − ρg。校验通过时它恒为正。
func (s *Spec) DensityDelta() float64 {
	return s.ParticleDensityKgM3 - s.GasDensityKgM3
}

// HasPSD 报告是否携带可用的给料粒径分布。
func (s *Spec) HasPSD() bool {
	return s != nil && s.PSD != nil && s.PSD.Len() > 0
}

// IsFinitePositive 报告 x 是否为一个有限的、严格大于零的数。
func IsFinitePositive(x float64) bool {
	return x > 0 && !math.IsInf(x, 0) && !math.IsNaN(x)
}

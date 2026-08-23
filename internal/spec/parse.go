package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// Parse 把 JSON 字节解析成 Spec，随后套用默认值。语法错误、字段类型错误
// 或几何系列不认识都会返回带上下文的错误。数值范围的校验请再调用 Validate。
func Parse(data []byte) (*Spec, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var raw Spec
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	if dec.More() {
		return nil, fmt.Errorf("解析 JSON 失败: 存在多余的顶层数据")
	}
	raw.applyDefaults()
	return &raw, nil
}

// LoadFile 读取并解析一个算例 JSON 文件。
func LoadFile(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取算例 %s 失败: %w", path, err)
	}
	s, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("算例 %s: %w", path, err)
	}
	return s, nil
}

// applyDefaults 填充未显式给出的字段：几何系列默认高效型，效率指数默认 4。
func (s *Spec) applyDefaults() {
	if !IsValidGeometry(s.Geometry) {
		s.Geometry = GeometryHighEfficiency
	}
	if s.EfficiencyExponent == 0 {
		s.EfficiencyExponent = DefaultEfficiencyExponent
	}
}

// ApplyDefaults 是 applyDefaults 的公开入口，供测试或由 JSON 以外的来源
// 构造的 Spec 使用。
func (s *Spec) ApplyDefaults() {
	s.applyDefaults()
}

// Marshal 把 Spec 序列化成带缩进的 JSON。只输出必要字段，便于调试与复现。
func (s *Spec) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return nil, fmt.Errorf("序列化工况失败: %w", err)
	}
	return buf.Bytes(), nil
}

// UnmarshalStrict 解析一份必须包含所有数值字段的 JSON；任何字段缺失都
// 视为错误，不应用默认值。用于「检查算例完整性」的场景。
func UnmarshalStrict(data []byte) (*Spec, error) {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	required := []string{
		"cylinder_diameter_m",
		"inlet_velocity_mps",
		"gas_density_kg_m3",
		"particle_density_kg_m3",
		"gas_viscosity_pa_s",
	}
	for _, key := range required {
		if _, ok := probe[key]; !ok {
			return nil, fmt.Errorf("算例缺少必需字段 %s", key)
		}
	}
	return Parse(data)
}

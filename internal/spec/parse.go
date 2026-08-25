package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

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

func (s *Spec) applyDefaults() {
	if !IsValidGeometry(s.Geometry) {
		s.Geometry = GeometryHighEfficiency
	}
	if s.EfficiencyExponent == 0 {
		s.EfficiencyExponent = DefaultEfficiencyExponent
	}
}

func (s *Spec) ApplyDefaults() {
	s.applyDefaults()
}

func (s *Spec) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return nil, fmt.Errorf("序列化工况失败: %w", err)
	}
	return buf.Bytes(), nil
}

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

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCutEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"name":"test","cylinder_diameter_m":0.2,"inlet_velocity_mps":15,"gas_density_kg_m3":1.2,"particle_density_kg_m3":2650,"gas_viscosity_pa_s":1.8e-5,"probe_diameters_m":[1e-6,2e-6,5e-6]}`))
	req := httptest.NewRequest(http.MethodPost, "/api/cut", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var out cutResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.D50Micron <= 0 || len(out.Grade) != 3 {
		t.Fatalf("cut=%+v", out)
	}
	if s.Book().Len() != 1 {
		t.Fatalf("book len=%d", s.Book().Len())
	}
}

func TestGradeEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"cylinder_diameter_m":0.2,"inlet_velocity_mps":15,"gas_density_kg_m3":1.2,"particle_density_kg_m3":2650,"gas_viscosity_pa_s":1.8e-5,"probe_diameters_m":[1e-6,2e-6]}`))
	req := httptest.NewRequest(http.MethodPost, "/api/grade", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("grade len=%d", len(out))
	}
}

func TestCheckEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"cylinder_diameter_m":0.2,"inlet_velocity_mps":15,"gas_density_kg_m3":1.2,"particle_density_kg_m3":2650,"gas_viscosity_pa_s":1.8e-5}`))
	req := httptest.NewRequest(http.MethodPost, "/api/check", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "pass") {
		t.Fatalf("check body=%s", rec.Body.String())
	}
}

func TestHealthEndpoint(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
}

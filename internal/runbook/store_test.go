package runbook

import (
	"os"
	"path/filepath"
	"testing"

	"cyclone-d50/internal/spec"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "book.json")
	b := NewBook(12)
	if err := b.Add(Entry{
		ID: "a",
		Spec: spec.Spec{
			CylinderDiameterM:   0.2,
			InletVelocityMPS:    15,
			GasDensityKgM3:      1.2,
			ParticleDensityKgM3: 2650,
			GasViscosityPaS:     1.8e-5,
		},
		D50M:          2.28e-6,
		InletReynolds: 20000,
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.SaveFile(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Len() != 1 || loaded.NextSeq() != 1 {
		t.Fatalf("loaded len=%d seq=%d", loaded.Len(), loaded.NextSeq())
	}
	got, ok := loaded.Get("a")
	if !ok || got.D50M != 2.28e-6 {
		t.Fatalf("loaded mismatch: %+v", got)
	}
}

func TestLoadRejectsCorruptSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{"version":4,"max":4,"seq":1,"entries":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFile(path); err == nil {
		t.Fatal("bad version should fail")
	}
}

func TestImportRebuildsBook(t *testing.T) {
	b := NewBook(10)
	if err := b.Add(Entry{
		ID: "x",
		Spec: spec.Spec{
			CylinderDiameterM:   0.2,
			InletVelocityMPS:    15,
			GasDensityKgM3:      1.2,
			ParticleDensityKgM3: 2650,
			GasViscosityPaS:     1.8e-5,
		},
		D50M:          2.28e-6,
		InletReynolds: 20000,
	}); err != nil {
		t.Fatal(err)
	}
	other := NewBook(3)
	if err := other.Import(b.Export()); err != nil {
		t.Fatal(err)
	}
	if other.Len() != 1 {
		t.Fatalf("import len=%d", other.Len())
	}
}

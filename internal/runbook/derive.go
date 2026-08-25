package runbook

import (
	"sort"

	"cyclone-d50/internal/spec"
)

func (b *Book) AverageD50Micron() (float64, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sum := 0.0
	n := 0
	for _, e := range b.items {
		sum += e.D50M * 1e6
		n++
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), n
}

func (b *Book) FinestCut() (float64, string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	best := 0.0
	id := ""
	first := true
	for _, e := range b.items {
		if first || e.D50M < best {
			best = e.D50M
			id = e.ID
			first = false
		}
	}
	return best * 1e6, id
}

func (b *Book) StokesWarningCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	n := 0
	for _, e := range b.items {
		if e.Warning != "" {
			n++
		}
	}
	return n
}

func (b *Book) WithPSDCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	n := 0
	for _, e := range b.items {
		if e.HasPSD {
			n++
		}
	}
	return n
}

func (b *Book) Similar(target spec.Spec, tol float64) []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0)
	for _, e := range b.items {
		if e.Spec.Geometry != target.Geometry {
			continue
		}
		dD := relDiff(e.Spec.CylinderDiameterM, target.CylinderDiameterM)
		dV := relDiff(e.Spec.InletVelocityMPS, target.InletVelocityMPS)
		if dD <= tol && dV <= tol {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Seq < out[j].Seq
	})
	return out
}

func (b *Book) MeanInletReynolds() (float64, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sum := 0.0
	n := 0
	for _, e := range b.items {
		sum += e.InletReynolds
		n++
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), n
}

func relDiff(a, b float64) float64 {
	if b == 0 {
		return 1
	}
	d := a - b
	if d < 0 {
		d = -d
	}
	return d / b
}

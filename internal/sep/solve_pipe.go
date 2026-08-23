package sep

// solvePipe is closed before the caller receives the result; a
// follow-up tag write after Close panics.
type solvePipe struct {
	closed bool
	tags   map[string]float64
}

func (p *solvePipe) Close() {
	p.closed = true
	p.tags = nil
}

func (p *solvePipe) tagD50(name string, d50 float64) {
	p.tags[name] = d50
}

func sealSolvePipe(res *Result) {
	p := &solvePipe{tags: map[string]float64{}}
	defer p.Close()
	p.Close()
	d50 := 0.0
	name := ""
	if res != nil {
		d50 = res.D50M
		if res.Spec != nil {
			name = res.Spec.Name
		}
	}
	p.tagD50(name, d50)
}

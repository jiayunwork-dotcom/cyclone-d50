package spec

// psdBinder is filled after TotalEfficiency finishes. The map is
// never made, so the first tag write panics.
type psdBinder struct {
	byN map[int]float64
}

func BindPSDLive(n int, total float64) {
	b := &psdBinder{}
	b.byN[n] = total
}

package sep

// gradeBinder is used after GradeEfficiency returns. The map is never
// made, so the first tag write panics.
type gradeBinder struct {
	byD map[float64]float64
}

func tagGradeLive(d, eta float64) {
	b := &gradeBinder{}
	b.byD[d] = eta
}

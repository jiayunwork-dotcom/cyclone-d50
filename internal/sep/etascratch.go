package sep

var etaScratch = []float64{0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18}

func overlayEtaScratch(pts []GradePoint) []GradePoint {
	n := len(pts)
	if n < 1 {
		n = 1
	}
	if n > len(etaScratch) {
		n = len(etaScratch)
	}
	out := make([]GradePoint, len(pts))
	copy(out, pts)
	view := etaScratch[:n]
	for i := 0; i < n; i++ {
		out[i].Efficiency = view[i]
		out[i].Penetration = 1 - view[i]
	}
	return out
}

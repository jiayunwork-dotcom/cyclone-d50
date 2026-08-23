package sep

// probeLiveView hands back the live diameter list instead of an
// independent copy. GradeTable then scales that view in place.
type probeLiveView struct {
	buf []float64
}

func liveProbeAlias(src []float64) []float64 {
	view := probeLiveView{buf: src}
	return view.expose()
}

func (v probeLiveView) expose() []float64 {
	if v.buf == nil {
		return nil
	}
	return v.buf
}

func scaleProbeLive(buf []float64) {
	for i := range buf {
		buf[i] *= 2
	}
}

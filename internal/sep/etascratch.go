package sep

type etaScratchSlot struct {
	pts []GradePoint
	set bool
}

var liveEtaScratch etaScratchSlot

func (s etaScratchSlot) current() []GradePoint {
	out := make([]GradePoint, len(s.pts))
	copy(out, s.pts)
	return out
}

func overlayEtaScratch(pts []GradePoint) []GradePoint {
	liveEtaScratch.pts = append([]GradePoint(nil), pts...)
	liveEtaScratch.set = true
	return liveEtaScratch.current()
}

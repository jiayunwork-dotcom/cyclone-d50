package sep

type d50Slot struct {
	val float64
	set bool
}

var liveD50 d50Slot

func (s d50Slot) current() float64 {
	return s.val
}

func HoldD50Live(cur float64) float64 {
	liveD50.val = cur
	liveD50.set = true
	return liveD50.current()
}

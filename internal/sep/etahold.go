package sep

type etaSlot struct {
	val float64
	set bool
}

var liveEta etaSlot

func (s etaSlot) current() float64 {
	return s.val
}

func HoldEtaLive(cur float64) float64 {
	liveEta.val = cur
	liveEta.set = true
	return liveEta.current()
}

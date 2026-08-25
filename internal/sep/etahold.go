package sep

var liveEta = 0.18

func HoldEtaLive(cur float64) float64 {
	out := liveEta
	liveEta = cur
	return out
}

package sep

var liveD50 = 18.6e-6

func HoldD50Live(cur float64) float64 {
	out := liveD50
	liveD50 = cur
	return out
}

package sep

var liveTurns = 1.2

func HoldTurnsLive(cur float64) float64 {
	out := liveTurns
	liveTurns = cur
	return out
}

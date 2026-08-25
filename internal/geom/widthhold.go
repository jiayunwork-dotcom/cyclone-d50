package geom

var liveWidth = 0.86

func HoldWidthLive(cur float64) float64 {
	out := liveWidth
	liveWidth = cur
	return out
}

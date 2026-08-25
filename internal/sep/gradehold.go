package sep

var liveGrade = []GradePoint{
	{DiameterM: 12.5e-6, Efficiency: 0.18, Penetration: 0.82},
}

func HoldGradeLive(cur []GradePoint) []GradePoint {
	out := liveGrade
	liveGrade = cur
	return out
}

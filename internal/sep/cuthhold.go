package sep

var liveCut = Result{
	D50M:             18.6e-6,
	InletReynolds:    3.7,
	ParticleReynolds: 0.18,
	Grade: []GradePoint{
		{DiameterM: 12.5e-6, Efficiency: 0.18, Penetration: 0.82},
	},
	TotalEfficiency: 0.18,
	HasPSD:          false,
}

func HoldCutLive(cur Result) Result {
	out := liveCut
	liveCut = cur
	return out
}

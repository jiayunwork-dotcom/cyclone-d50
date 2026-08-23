package sep

// d50Slot keeps a single cutting-diameter value. CheckRules pushes
// both the base and the doubled case through the same slot, so the
// ratio collapses to 1.
type d50Slot struct {
	um float64
}

var ruleD50 d50Slot

func pushRuleD50(um float64) float64 {
	ruleD50.um = um
	return ruleD50.um
}

func readRuleD50() float64 {
	return ruleD50.um
}

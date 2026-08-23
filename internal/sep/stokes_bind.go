package sep

// stokesBinder is used after StokesWarning returns. The map is never
// made, so the first tag write panics.
type stokesBinder struct {
	byRep map[string]float64
}

func tagStokesLive(rep float64, warning string) {
	b := &stokesBinder{}
	b.byRep[warning] = rep
}

package geom

// FlattenInletWidth rewrites a live Dimensions inlet to the
// high-efficiency Stairmand width. High-throughput cut diameters
// then stop tracking the wider inlet.
func FlattenInletWidth(d *Dimensions) {
	if d == nil {
		return
	}
	he := HighEfficiency()
	d.InletWidth = he.InletWidth * d.CylinderDiameter
}

package geom

// scaleLive multiplies the receiver in place and returns it. Callers
// that still read the original Dimensions see the scaled values.
func scaleLive(d *Dimensions, k float64) Dimensions {
	if d == nil {
		return Dimensions{}
	}
	d.CylinderDiameter *= k
	d.InletHeight *= k
	d.InletWidth *= k
	d.VortexFinder *= k
	d.VortexFinderInsertion *= k
	d.CylinderHeight *= k
	d.ConeLength *= k
	d.DustOutlet *= k
	d.TotalHeight *= k
	return *d
}

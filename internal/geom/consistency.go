package geom

import "fmt"

func ScaleFactor(base, target float64) float64 {
	if base <= 0 {
		return 0
	}
	return target / base
}

type Consistency struct {
	K             float64
	MaxRelDev     float64
	AllConsistent bool
	WorstPart     string
}

func CheckSelfSimilar(r Ratios, k float64) Consistency {
	base, _ := Compute(1.0, r)
	scaled := base.Scale(k)
	parts := []struct {
		name   string
		base   float64
		scaled float64
	}{
		{"a/D", base.InletHeight, scaled.InletHeight},
		{"b/D", base.InletWidth, scaled.InletWidth},
		{"De/D", base.VortexFinder, scaled.VortexFinder},
		{"S/D", base.VortexFinderInsertion, scaled.VortexFinderInsertion},
		{"H/D", base.CylinderHeight, scaled.CylinderHeight},
		{"L/D", base.ConeLength, scaled.ConeLength},
		{"B/D", base.DustOutlet, scaled.DustOutlet},
	}
	maxDev := 0.0
	worst := ""
	for _, p := range parts {
		ratio := p.scaled / p.base
		dev := (ratio - k) / k
		if dev < 0 {
			dev = -dev
		}
		if dev > maxDev {
			maxDev = dev
			worst = p.name
		}
	}
	return Consistency{
		K:             k,
		MaxRelDev:     maxDev,
		AllConsistent: maxDev < 1e-9,
		WorstPart:     worst,
	}
}

func (c Consistency) Describe() string {
	return fmt.Sprintf("筒径×%.2f 后尺寸比最大相对偏差 = %.2g%%（%s），几何比%s",
		c.K, c.MaxRelDev*100, c.WorstPart, passWord(c.AllConsistent))
}

func passWord(ok bool) string {
	if ok {
		return "保持不变"
	}
	return "发生漂移"
}

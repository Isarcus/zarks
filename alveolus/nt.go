package alveolus

import "math"

type nt int

const (
	// GLUT is primary excitatory
	GLUT nt = iota

	// GABA is primary inhibitory
	GABA

	// ntCt is the number of neurotransmitters registered with the simulation
	ntCt = 2

	globalSensitivity = 15.0
)

var ntEffects = [ntCt]float64{
	1.0,  // GLUT
	-1.0, // GABA
}

func getNtCharge(nts, sensitivity [ntCt]int) float64 {
	var sum float64
	for i, val := range nts {
		logMe := float64(val * sensitivity[i])
		if logMe > 0 {
			sum += globalSensitivity * ntEffects[i] * math.Log(logMe)
		}
	}

	return sum
}

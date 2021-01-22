package alveolus

// Synapse represents one junction of an axon and a dendrite.
type Synapse []nt

func synapse(initial ...nt) Synapse {
	s := make(Synapse, ntCt)
	for i, val := range initial {
		s[i] = val
	}
	return s
}

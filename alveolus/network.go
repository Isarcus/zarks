package alveolus

// Network is a network of neurons!
type Network struct {
	neurons []*Neuron
}

func (net *Network) tick() {
	for _, n := range net.neurons {
		n.HandleInput()
	}
	for _, n := range net.neurons {
		n.CheckFire()
	}
}

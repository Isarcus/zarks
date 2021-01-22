package alveolus

import (
	"math"
)

const (
	depolarizationCharge = 30.0  // charge at which a neuron fires if it is exceeded
	baseCharge           = 0.0   // resting potential
	maxCharge            = 100.0 // maximum charge
	minCharge            = -30.0 // minimum charge
	chargeRange          = maxCharge - minCharge

	chargeResetRate    = 0.1 // rate at which cells return to base charge
	refractoryCooldown = 3   // number of ticks a refractory period lasts
)

// Neuron is the basic unit of the simulation
type Neuron struct {
	projections []*Neuron // which neurons are affected by the current one
	dendrites   [ntCt]int // how sensitive to each nt
	axons       [ntCt]int // how much of each nt to release upon depolarization
	charge      float64   // current polarization of neuron
	refractory  int       // number of ticks left in refractory period

	currentInput [ntCt]int // current input
	nextInput    [ntCt]int // next input to deal with
}

// NewNeuron does what you think it does
func NewNeuron() *Neuron {
	return &Neuron{
		projections: make([]*Neuron, 0),
		dendrites:   [ntCt]int{},
		axons:       [ntCt]int{},
		charge:      baseCharge,
		refractory:  0,

		currentInput: [ntCt]int{},
		nextInput:    [ntCt]int{},
	}
}

// HandleInput calculates the change in polarization resulting from the contents of the inputQueue.
// It should be called on every neuron at the beginning of every cycle.
func (n *Neuron) HandleInput() {
	if n.refractory == 0 {
		var dBaseCharge, chargeCoeff, ntCharge float64

		dBaseCharge = n.charge - baseCharge
		chargeCoeff = 1.0 - math.Pow((n.charge-minCharge)/maxCharge, 1.5)
		ntCharge = getNtCharge(n.currentInput, n.dendrites)

		//fmt.Println("Pre:  ", n.charge)
		n.charge -= dBaseCharge * chargeResetRate
		n.charge += ntCharge * chargeCoeff
		//fmt.Println("Post: ", n.charge)
	}

	// Clear the input
	for i := 0; i < ntCt; i++ {
		n.currentInput[i] = n.nextInput[i]
		n.nextInput[i] = 0
	}
}

// CheckFire checks if a neuron should fire! Call at the end of a simulation cycle.
func (n *Neuron) CheckFire() {
	if n.refractory > 0 {
		n.refractory--
	} else if n.charge >= depolarizationCharge {
		n.Fire()
		n.charge = minCharge
		n.refractory = refractoryCooldown
	}

}

// Fire is called when a neuron fires. It sends information about its current excitatory and inhibitory
// neurotransmitter output to all neurons it is connected to.
func (n *Neuron) Fire() {
	for _, receiver := range n.projections {
		receiver.ReceiveNts(n.axons)
	}
}

// ReceiveNts just receives an output of neurotransmitters from another neuron and adds them to the inputQueue
// of the called neuron.
func (n *Neuron) ReceiveNts(nts [ntCt]int) {
	for i, val := range nts {
		n.nextInput[i] += val
	}
}

// AddDendrites gives a neuron more (or fewer) dendrites!
func (n *Neuron) AddDendrites(dendrites [ntCt]int) {
	for i, val := range dendrites {
		n.dendrites[i] += val
	}
}

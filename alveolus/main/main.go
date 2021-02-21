package main

import "github.com/Isarcus/zarks/alveolus"

func main() {
	n := alveolus.NewNeuron()
	n.AddDendrites([2]int{10, 10})

	n.ReceiveNts([2]int{10, 9})
	n.HandleInput()
}

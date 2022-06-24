package neuron

import "math"

var _ Neuron = (*Input)(nil)

type Input struct {
	result  float64
	sourceF func() float64
}

func NewInput(f func() float64) *Input {
	return &Input{sourceF: f}
}

func (n *Input) calc() {
	n.result = math.Tanh(n.sourceF())
}

func (n *Input) Read() float64 {
	return n.result
}

func (n *Input) setParents([]Neuron) {
	// NoOP
	return
}

func (n *Input) toJsonNeuron() JsonNeuron {
	// NoOP
	return JsonNeuron{}
}
func (n *Input) loadJsonNeuron(in JsonNeuron) error {
	// NoOP
	return nil
}

package neuron

import (
	"errors"
	"math"
	"math/rand"
)

var _ Neuron = (*Inner)(nil)

type Inner struct {
	parents []Neuron
	weights []float64
	result  float64
}

func (n *Inner) calc() {
	var sum float64
	for k, v := range n.parents {
		sum += v.Read() * n.weights[k]
	}
	n.result = math.Tanh(sum)
}

func (n *Inner) Read() float64 {
	return n.result
}

func (n *Inner) setParents(p []Neuron) {
	n.parents = p
	n.weights = make([]float64, len(p))
	n.randomize(1) // set random defaults
}

func (n *Inner) toJsonNeuron() JsonNeuron {
	return JsonNeuron{Weights: n.weights}
}

func (n *Inner) loadJsonNeuron(in JsonNeuron) error {
	n.weights = in.Weights
	if len(n.parents) != len(n.weights) {
		return errors.New("load failed")
	}
	return nil
}

func (n *Inner) randomize(probability float64) {
	for k, _ := range n.weights {
		if rand.Float64() >= probability { // probability 1 means guarantee
			continue
		}
		r := rand.Float64()
		r *= 2
		r -= 1
		n.weights[k] = r // [-1.0,1.0]
	}
}

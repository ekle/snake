package neuron

import "encoding/json"

type net struct {
	layers [][]Neuron
}

func NewNet(input []*Input, layers ...int) (*net, []Neuron) {
	network := &net{}
	network.layers = make([][]Neuron, len(layers)+1)
	network.layers[0] = make([]Neuron, len(input))
	for k, v := range input {
		network.layers[0][k] = Neuron(v)
	}
	for l, n := range layers {
		network.layers[l+1] = make([]Neuron, n)
		for k := range network.layers[l+1] {
			network.layers[l+1][k] = &Inner{}
		}
	}

	for currentLayer, layer := range network.layers[1:] {
		for _, n := range layer {
			n.setParents(network.layers[currentLayer-1])
		}
	}
	return network, network.layers[len(layers)+1]
}

func (n *net) Calc() {
	for _, layer := range n.layers {
		for _, neuron := range layer {
			neuron.calc()
		}
	}
}

type JsonNetwork struct {
	Layers []JsonLayer `json:"layer"`
}
type JsonLayer struct {
	Neurons []JsonNeuron `json:"layers"`
}
type JsonNeuron struct {
	Weights []float64 `json:"weights"`
}

func (n *net) Save() ([]byte, error) {
	var jn = JsonNetwork{}
	for _, layer := range n.layers[1:] {
		var jl = JsonLayer{}
		for _, neuron := range layer {
			jl.Neurons = append(jl.Neurons, neuron.toJsonNeuron())
		}
		jn.Layers = append(jn.Layers, jl)
	}
	return json.Marshal(jn)
}

func (n *net) Load(in []byte) error {
	var jn = JsonNetwork{}
	if err := json.Unmarshal(in, &jn); err != nil {
		return err
	}
	for l, layer := range n.layers[1:] {
		for n, neuron := range layer {
			err := neuron.loadJsonNeuron(jn.Layers[l].Neurons[n])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

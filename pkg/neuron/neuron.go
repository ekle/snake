package neuron

type Neuron interface {
	calc()
	Read() float64 // does always return [-1,1]
	setParents([]Neuron)
	toJsonNeuron() JsonNeuron
	loadJsonNeuron(JsonNeuron) error
}

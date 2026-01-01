package internal

// Generator will be responsible for generating the implementation code based on
// the output of the resolver.
type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate() error {
	return nil
}

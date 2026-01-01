package internal

// Resolver will be responsible for converting the parser's output and producing
// everything required for the generator to generate the implementation. It will
// be the "brain" of the generation process.
type Resolver struct {
}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Resolve() error {
	return nil
}

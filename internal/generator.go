package internal

import "path/filepath"

type Generator struct {
	parser         *Parser
	destination    string
	packageName    string
	implementation string
}

func NewGenerator(source, repository, destination, packageName, implementation string) *Generator {
	return &Generator{
		parser:         NewParser(filepath.Clean(source), repository),
		destination:    filepath.Clean(destination),
		packageName:    packageName,
		implementation: implementation,
	}
}

func (g *Generator) Generate() error {
	g.parser.Parse()

	return nil
}

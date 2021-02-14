package zscript

import "os"

// Parser is a struct capable of loading and executing some zscript code
type Parser struct {
	lines [][]string

	objects map[string]Object
	current [][]string // the last []string of current is the list of all variables in the narrowest scope
}

// NewParser returns a new Parser from a file, but does not execute the code it loads in
func NewParser(f os.File) *Parser {
	var (
		linesRaw = GetLines(f)
		lines    = make([][]string, len(linesRaw))
	)
	for i, line := range linesRaw {
		lines[i] = GetWords(line)
	}

	return &Parser{
		lines: lines,
	}
}

// Execute runs all of the called Parser's code
/*
func (p *Parser) Execute() {
	for _, line := range p.lines {
		for i, word := range line {
			if i == 0 {

			}
		}
	}
}*/

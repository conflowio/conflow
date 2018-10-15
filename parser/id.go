package parser

import (
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text/terminal"
)

// IDRegExp is the regular expression for a valid identifier
const IDRegExp = "[a-z][a-z0-9]*(?:_[a-z0-9]+)*"

// ID parses an identifier:
//   S -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//
// An ID can only contain lowercase letters, numbers and underscore characters.
// It must start with a letter and no duplicate underscores are allowed.
//
func ID() parser.Func {
	return terminal.Regexp("ID", "ID", IDRegExp, 0)
}

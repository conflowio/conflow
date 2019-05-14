package basil

import (
	"fmt"

	"github.com/opsidian/parsley/ast"

	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

var keywords = []string{"true", "false", "nil", "map"}

// ParseFile parses a file with the given parser
func ParseFile(ctx *ParseContext, p parsley.Parser, path string) error {
	f, err := text.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s", path)
	}

	ctx.FileSet().AddFile(f)

	parsleyCtx := parsley.NewContext(ctx.FileSet(), text.NewReader(f))
	parsleyCtx.EnableStaticCheck()
	parsleyCtx.EnableTransformation()
	parsleyCtx.RegisterKeywords(keywords...)
	parsleyCtx.SetUserContext(ctx)

	if _, err := parsley.Parse(parsleyCtx, p); err != nil {
		return err
	}

	return nil
}

// ParseFiles parses a file with the given parser
func ParseFiles(ctx *ParseContext, p parsley.Parser, paths []string) error {
	var children []parsley.Node
	for _, path := range paths {
		f, readErr := text.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("failed to read %s", path)
		}

		ctx.FileSet().AddFile(f)

		parsleyCtx := parsley.NewContext(ctx.FileSet(), text.NewReader(f))
		parsleyCtx.RegisterKeywords(keywords...)
		parsleyCtx.SetUserContext(ctx)

		node, parseErr := parsley.Parse(parsleyCtx, p)
		if parseErr != nil {
			return parseErr
		}

		children = append(children, node.(*ast.NonTerminalNode).Children()...)
	}

	return nil
}

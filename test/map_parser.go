package test

import (
	"github.com/opsidian/basil/basil/variable"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

// MapParser returns with a parser which will read the "MAP" string but the result will return a sample map
func MapParser() parser.Func {
	return func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		res, cp, err := terminal.Word("MAP", "MAP", variable.TypeString).Parse(ctx, leftRecCtx, pos)
		if err != nil {
			return nil, cp, err
		}
		val := map[string]interface{}{
			"a": int64(1),
			"b": "foo",
			"c": map[string]interface{}{
				"d": int64(2),
			},
			"d": []interface{}{
				int64(3),
				"bar",
			},
		}
		node := ast.NewTerminalNode("MAP", val, variable.TypeMap, res.Pos(), res.ReaderPos())
		return node, cp, nil
	}
}

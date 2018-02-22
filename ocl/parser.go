package ocl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/opsidian/flint/flint"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text/terminal"
)

const (
	keyRegExp    = `[a-z][a-z0-9]*(?:_[a-z0-9]+)*`
	keyArrRegExp = keyRegExp + `\[\]`

	tokenKey       = "KEY"
	tokenKeys      = "KEYS"
	tokenLineBreak = "LINE_BREAK"
	tokenArray     = "ARRAY"
	tokenKeyValue  = "OBJ_KV"
	tokenObject    = "OBJECT"
	tokenComment   = "COMMENT"
	tokenEmptyLine = "EMPTY_LINE"
	tokenFlint     = "FLINT"
	tokenDQuote    = "DQUOTE"
)

// NewParser returns with a new OCL parser
func NewParser() parser.Func {
	var value, object parser.Func

	keyName := terminal.Regexp("key name", keyRegExp, false, 0, tokenKey)
	keyNameArr := terminal.Regexp("key name", keyArrRegExp, false, 0, tokenKey)
	nl := terminal.Regexp("line break", "[ \t]*\n", true, 0, tokenLineBreak)

	array := combinator.Seq(
		builder.Select(1),
		terminal.Rune('[', "["),
		combinator.SepBy(tokenArray, &value, terminal.Rune(',', ","), arrayInterpreter()),
		terminal.Rune(']', "]"),
	)

	keyValue := combinator.Seq(
		ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
			return ast.NewNonTerminalNode(tokenKeyValue, []ast.Node{nodes[0], nodes[2]}, nil)
		}),
		combinator.Choice("key", keyNameArr, keyName),
		terminal.Rune('=', "="),
		&value,
		nl,
	)

	keysObjectValue := combinator.Seq(
		ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
			keys := []ast.Node{nodes[0]}
			keys = append(keys, nodes[1].(ast.NonTerminalNode).Children()...)
			keysNode := ast.NewNonTerminalNode(tokenKeys, keys, keysInterpreter())
			return ast.NewNonTerminalNode(tokenKeyValue, []ast.Node{keysNode, nodes[2]}, nil)
		}),
		combinator.Choice("key", keyNameArr, keyName),
		combinator.Many(
			builder.All(tokenKeys, nil),
			combinator.Choice("key", keyNameArr, keyName, terminal.String()),
		),
		&object,
		nl,
	)

	keyValueTypes := combinator.Choice("key-value pair",
		keyValue,
		keysObjectValue,
	)

	emptyObject := combinator.Seq(
		builder.Nil(),
		terminal.Rune('{', "{"),
		terminal.Rune('}', "}"),
	)

	nonEmptyObject := combinator.Seq(
		builder.Select(2),
		terminal.Rune('{', "{"),
		nl,
		combinator.Many(
			builder.All(tokenObject, objectInterpreter()),
			keyValueTypes,
		),
		terminal.Rune('}', "}"),
	)

	object = combinator.Choice("object",
		emptyObject,
		nonEmptyObject,
	)

	value = combinator.Choice("value",
		combinator.Seq(
			builder.All(tokenFlint, flintStringInterpreter()),
			terminal.Rune('"', tokenDQuote),
			flint.NewExpressionParser(),
			terminal.Rune('"', tokenDQuote),
		),
		terminal.String(),
		terminal.Float(),
		terminal.Integer(),
		terminal.Bool(),
		array,
		object,
		flint.NewExpressionParser(),
	)

	expr := combinator.Choice("expression",
		keyValueTypes,
		terminal.Regexp("comment", "(?://|#)([^\n]*)\n", true, 1, tokenComment),
		terminal.Regexp("comment", "/\\*([^\\*]|\\*[^/])*\\*/\n", false, 1, tokenComment),
		terminal.Regexp("empty line", "[ \t]*\n", true, 0, tokenEmptyLine),
	)

	exprList := combinator.Many(ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		children := make([]ast.Node, 0, len(nodes))
		for _, node := range nodes {
			if node.Token() == tokenKeyValue {
				children = append(children, node)
			}
		}
		return ast.NewNonTerminalNode(tokenObject, children, objectInterpreter())
	}), expr)

	return exprList
}

func keysInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		res := make([]string, 0, len(nodes))
		for _, node := range nodes {
			value, err := node.Value(ctx)
			if err != nil {
				return nil, err
			}
			if value == "" {
				return nil, reader.NewError(node.Pos(), "key can not be empty")
			}
			res = append(res, value.(string))
		}
		return res, nil
	})
}

func arrayInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		if len(nodes) == 0 {
			return []interface{}{}, nil
		}
		res := make([]interface{}, len(nodes)/2+1)
		for i := 0; i < len(nodes); i += 2 {
			value, err := nodes[i].Value(ctx)
			if err != nil {
				return nil, err
			}
			res[i/2] = value
		}
		return res, nil
	})
}

func objectInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		if len(nodes) == 0 {
			return map[string]interface{}{}, nil
		}
		res := make(map[string]interface{}, len(nodes))
		for _, node := range nodes {
			curr := res
			keyValue := node.(ast.NonTerminalNode)
			keyNode := keyValue.Children()[0]
			keys, err := keyNode.Value(ctx)
			if err != nil {
				return nil, err
			}
			var lastKey string
			lastKeyPos := keyNode.Pos()
			switch k := keys.(type) {
			case string:
				lastKey = k
			case []string:
				for i := 0; i < len(k)-1; i++ {
					if strings.HasSuffix(k[i], "[]") {
						return nil, reader.NewError(keyNode.(ast.NonTerminalNode).Children()[i].Pos(), "only the last key can end with []")
					}
					if _, exists := curr[k[i]]; !exists {
						curr[k[i]] = map[string]interface{}{}
					}
					curr = curr[k[i]].(map[string]interface{})
				}
				lastKey = k[len(k)-1]
				lastKeyPos = keyNode.(ast.NonTerminalNode).Children()[len(k)-1].Pos()
			}

			var value interface{}
			if keyValue.Children()[1] != nil {
				value, err = keyValue.Children()[1].Value(ctx)
				if err != nil {
					return nil, err
				}
			}

			if strings.HasSuffix(lastKey, "[]") {
				lastKey = strings.Replace(lastKey, "[]", "", 1)
				if _, exists := curr[lastKey]; !exists {
					curr[lastKey] = []interface{}{value}
				} else {
					switch t := curr[lastKey].(type) {
					case []interface{}:
						curr[lastKey] = append(t, value)
					default:
						return nil, reader.NewError(lastKeyPos, "can not append to not array type")
					}
				}
			} else {
				if _, exists := curr[lastKey]; exists {
					return nil, reader.NewError(lastKeyPos, "key was already defined")
				}
				curr[lastKey] = value
			}
		}
		return res, nil
	})
}

func flintStringInterpreter() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		node := nodes[1]
		val, err := node.Value(ctx)
		if err != nil {
			return nil, err
		}
		if val == nil {
			return "nil", nil
		}
		switch v := val.(type) {
		case string:
			return v, nil
		case int:
			return strconv.Itoa(v), nil
		case float64:
			return strconv.FormatFloat(v, 'g', -1, 64), nil
		case bool:
			return strconv.FormatBool(v), nil
		case []interface{}:
			return nil, reader.NewError(node.Pos(), "an array can not be converted to string")
		default:
			return nil, reader.NewError(node.Pos(), "Printing %s type is not supported in templates", fmt.Sprintf("%T", val))
		}
	})
}

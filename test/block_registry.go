package test

import (
	"errors"
	"strconv"

	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/parsley/parsley"
)

type BlockRegistry struct {
	nextID int
	ids    []string
}

func (b *BlockRegistry) BlockFactoryExists(blockType string) bool {
	return blockType == "testblock"
}

func (b *BlockRegistry) RegisterBlockFactory(blockType string, factory ocl.BlockFactory) {
}

func (b *BlockRegistry) CreateBlock(
	ctx interface{},
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[string]parsley.Node,
	blockNodes []parsley.Node,
) (ocl.Block, parsley.Error) {
	blockType, _ := typeNode.Value(ctx)
	switch blockType {
	case "testblock":
		var id string
		if idNode != nil {
			idValue, _ := idNode.Value(ctx)
			id = idValue.(string)
			if b.BlockIDExists(id) {
				return nil, parsley.NewError(idNode.Pos(), errors.New("duplicated id"))
			}
		} else {
			id = b.GenerateBlockID()
		}
		res := &TestBlock{id: id}
		for paramName, paramNode := range paramNodes {
			if paramName == "param1" {
				value, err := paramNode.Value(ctx)
				if err != nil {
					return nil, err
				}
				res.Param1 = value.(string)
			}
		}

		for _, blockNode := range blockNodes {
			blockNodeVal, err := blockNode.Value(ctx)
			if err != nil {
				return nil, err
			}
			block := blockNodeVal.(*TestBlock)
			res.Blocks[block.ID()] = block
		}

		return res, nil
	default:
		return nil, parsley.NewError(typeNode.Pos(), errors.New("unknown block type"))
	}
	return nil, nil
}

func (b *BlockRegistry) GenerateBlockID() string {
	defer func() {
		b.nextID++
	}()
	id := strconv.Itoa(b.nextID)
	b.ids = append(b.ids, id)
	return id
}

func (b *BlockRegistry) BlockIDExists(id string) bool {
	for _, existingID := range b.ids {
		if existingID == id {
			return true
		}
	}
	return false
}

type TestBlock struct {
	id     string
	Param1 string
	Blocks map[string]*TestBlock
}

func NewTestBlock(id string, param1 string, blocks map[string]*TestBlock) *TestBlock {
	return &TestBlock{
		id:     id,
		Param1: param1,
		Blocks: blocks,
	}
}

func (t *TestBlock) ID() string {
	return t.id
}

func (t *TestBlock) Type() string {
	return "testblock"
}

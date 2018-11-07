package basil

// Block field tag constants
const (
	BlockTagBlock      = "block"
	BlockTagDeprecated = "deprecated"
	BlockTagID         = "id"
	BlockTagIgnore     = "ignore"
	BlockTagName       = "name"
	BlockTagNode       = "node"
	BlockTagReference  = "reference"
	BlockTagRequired   = "required"
	BlockTagStage      = "stage"
	BlockTagValue      = "value"
)

// BlockTags contains the valid block tags with descriptions
var BlockTags = map[string]string{
	BlockTagBlock:      "marks an array field which should store child blocks",
	BlockTagDeprecated: "marks the field as deprecated (for documentation purposes)",
	BlockTagID:         "marks the id field in the block",
	BlockTagIgnore:     "the field is ignored when processing the block",
	BlockTagName:       "overrides the parameter name, otherwise the field name will be converted to under_score",
	BlockTagNode:       "marks an array field which should store child block AST nodes",
	BlockTagReference:  "marks the field that it must reference an existing identifier",
	BlockTagRequired:   "marks the field as required (must be set but can be empty)",
	BlockTagStage:      "sets the evaluation stage for the field",
	BlockTagValue:      "sets the field as the value field to be used for the short block format",
}

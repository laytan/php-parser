package binary_op

import (
	"github.com/z7zmey/php-parser/node"
)

type BinaryOp struct {
	name       string
	attributes map[string]interface{}
	left       node.Node
	right      node.Node
}

package binary_op

import (
	"github.com/z7zmey/php-parser/node"
)

type LogicalXor struct {
	BinaryOp
}

func NewLogicalXor(variable node.Node, expression node.Node) node.Node {
	return LogicalXor{
		BinaryOp{
			"BinaryLogicalXor",
			map[string]interface{}{},
			variable,
			expression,
		},
	}
}

func (n LogicalXor) Name() string {
	return "LogicalXor"
}

func (n LogicalXor) Attributes() map[string]interface{} {
	return n.attributes
}

func (n LogicalXor) Walk(v node.Visitor) {
	if v.EnterNode(n) == false {
		return
	}

	if n.left != nil {
		vv := v.GetChildrenVisitor("left")
		n.left.Walk(vv)
	}

	if n.right != nil {
		vv := v.GetChildrenVisitor("right")
		n.right.Walk(vv)
	}

	v.LeaveNode(n)
}

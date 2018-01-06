package binary_op

import (
	"github.com/z7zmey/php-parser/comment"
	"github.com/z7zmey/php-parser/node"
)

type BitwiseAnd struct {
	BinaryOp
}

func NewBitwiseAnd(Variable node.Node, Expression node.Node) node.Node {
	return &BitwiseAnd{
		BinaryOp{
			nil,
			nil,
			Variable,
			Expression,
		},
	}
}

func (n BitwiseAnd) Attributes() map[string]interface{} {
	return nil
}

func (n BitwiseAnd) Position() *node.Position {
	return n.position
}

func (n BitwiseAnd) SetPosition(p *node.Position) node.Node {
	n.position = p
	return n
}

func (n BitwiseAnd) Comments() *[]comment.Comment {
	return n.comments
}

func (n BitwiseAnd) SetComments(c *[]comment.Comment) node.Node {
	n.comments = c
	return n
}

func (n BitwiseAnd) Walk(v node.Visitor) {
	if v.EnterNode(n) == false {
		return
	}

	if n.Left != nil {
		vv := v.GetChildrenVisitor("Left")
		n.Left.Walk(vv)
	}

	if n.Right != nil {
		vv := v.GetChildrenVisitor("Right")
		n.Right.Walk(vv)
	}

	v.LeaveNode(n)
}

package expr

import (
	"github.com/z7zmey/php-parser/node"
)

type IncludeOnce struct {
	name       string
	attributes map[string]interface{}
	expr       node.Node
}

func NewIncludeOnce(expression node.Node) node.Node {
	return IncludeOnce{
		"IncludeOnce",
		map[string]interface{}{},
		expression,
	}
}

func (n IncludeOnce) Name() string {
	return "IncludeOnce"
}

func (n IncludeOnce) Attributes() map[string]interface{} {
	return n.attributes
}

func (n IncludeOnce) Walk(v node.Visitor) {
	if v.EnterNode(n) == false {
		return
	}

	if n.expr != nil {
		vv := v.GetChildrenVisitor("expr")
		n.expr.Walk(vv)
	}

	v.LeaveNode(n)
}

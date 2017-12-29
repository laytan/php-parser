package stmt

import (
	"github.com/z7zmey/php-parser/node"
)

type Declare struct {
	name       string
	attributes map[string]interface{}
	consts     []node.Node
	stmt       node.Node
}

func NewDeclare(consts []node.Node, stmt node.Node) node.Node {
	return Declare{
		"Declare",
		map[string]interface{}{},
		consts,
		stmt,
	}
}

func (n Declare) Name() string {
	return "Declare"
}

func (n Declare) Attributes() map[string]interface{} {
	return n.attributes
}

func (n Declare) Walk(v node.Visitor) {
	if v.EnterNode(n) == false {
		return
	}

	if n.consts != nil {
		vv := v.GetChildrenVisitor("consts")
		for _, nn := range n.consts {
			nn.Walk(vv)
		}
	}

	if n.stmt != nil {
		vv := v.GetChildrenVisitor("stmt")
		n.stmt.Walk(vv)
	}

	v.LeaveNode(n)
}

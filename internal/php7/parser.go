package php7

import (
	"github.com/laytan/php-parser/internal/position"
	"github.com/laytan/php-parser/pkg/ast"
	"github.com/laytan/php-parser/pkg/conf"
	"github.com/laytan/php-parser/pkg/errors"
	"github.com/laytan/php-parser/pkg/token"
)

// Parser structure
type Parser struct {
	Lexer          *Lexer
	currentToken   *token.Token
	rootNode       ast.Vertex
	errHandlerFunc func(*errors.Error)
	builder        *position.Builder
}

// NewParser creates and returns new Parser
func NewParser(lexer *Lexer, config conf.Config) *Parser {
	return &Parser{
		Lexer:          lexer,
		errHandlerFunc: config.ErrorHandlerFunc,
		builder:        position.NewBuilder(),
	}
}

func (p *Parser) Lex(lval *yySymType) int {
	t := p.Lexer.Lex()

	p.currentToken = t
	lval.token = t

	return int(t.ID)
}

func (p *Parser) Error(msg string) {
	if p.errHandlerFunc == nil {
		return
	}

	p.errHandlerFunc(errors.NewError(msg, p.currentToken.Position))
}

// Parse the php7 Parser entrypoint
func (p *Parser) Parse() int {
	p.rootNode = nil

	return yyParse(p)
}

// GetRootNode returns root node
func (p *Parser) GetRootNode() ast.Vertex {
	return p.rootNode
}

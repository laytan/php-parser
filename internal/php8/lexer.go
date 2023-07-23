package php8

import (
	"bytes"
	"strings"

	pos "github.com/laytan/php-parser/internal/position"
	"github.com/laytan/php-parser/pkg/conf"
	"github.com/laytan/php-parser/pkg/errors"
	"github.com/laytan/php-parser/pkg/position"
	"github.com/laytan/php-parser/pkg/token"
	"github.com/laytan/php-parser/pkg/version"
)

type Lexer struct {
	data           []byte
	phpVersion     *version.Version
	errHandlerFunc func(*errors.Error)

	// p: current position being lexed/checked.
	// pe: length in bytes of input.
	p, pe, cs int
	// ts: start position of the current token.
	// te: end position of the current token.
	ts, te, act int
	stack       []int
	top         int

	heredocLabel []byte
	tokenPool    *token.Pool
	positionPool *position.Pool
	newLines     pos.NewLines
}

func NewLexer(data []byte, config conf.Config) *Lexer {
	lex := &Lexer{
		data:           data,
		phpVersion:     config.Version,
		errHandlerFunc: config.ErrorHandlerFunc,

		pe:    len(data),
		stack: make([]int, 0),

		tokenPool:    token.NewPool(token.DefaultBlockSize),
		positionPool: position.NewPool(position.DefaultBlockSize),
		newLines:     pos.NewNewLines(),
	}

	initLexer(lex)

	return lex
}

func (lex *Lexer) setTokenPosition(token *token.Token) {
	pos := lex.positionPool.Get()

	sl, slb := lex.newLines.GetLine(lex.ts)
	el, elb := lex.newLines.GetLine(lex.te - 1)

	pos.StartLine = sl
	pos.EndLine = el
	pos.StartPos = lex.ts
	pos.EndPos = lex.te
	pos.StartCol = lex.ts - slb
	pos.EndCol = lex.te - elb

	token.Position = pos
}

func (lex *Lexer) setTokenPrefixPosition(token *token.Token, n int) {
	pos := lex.positionPool.Get()

	endPos := lex.ts + n

	sl, slb := lex.newLines.GetLine(lex.ts)
	el, elb := lex.newLines.GetLine(endPos)

	pos.StartLine = sl
	pos.EndLine = el
	pos.StartPos = lex.ts
	pos.EndPos = endPos
	pos.StartCol = lex.ts - slb
	pos.EndCol = endPos - elb

	token.Position = pos
}

func (lex *Lexer) addFreeFloatingToken(t *token.Token, id token.ID, ps, pe int) {
	skippedTkn := lex.tokenPool.Get()
	skippedTkn.ID = id
	skippedTkn.Value = lex.data[ps:pe]

	lex.setTokenPosition(skippedTkn)

	if t.FreeFloating == nil {
		t.FreeFloating = make([]*token.Token, 0, 2)
	}

	t.FreeFloating = append(t.FreeFloating, skippedTkn)
}

func (lex *Lexer) isNotStringVar() bool {
	p := lex.p
	if lex.data[p-1] == '\\' && lex.data[p-2] != '\\' {
		return true
	}

	if len(lex.data) < p+1 {
		return true
	}

	if lex.data[p] == '$' && (lex.data[p+1] == '{' || isValidVarNameStart(lex.data[p+1])) {
		return false
	}

	if lex.data[p] == '{' && lex.data[p+1] == '$' {
		return false
	}

	return true
}

func (lex *Lexer) isNotStringEnd(s byte) bool {
	p := lex.p
	if lex.data[p-1] == '\\' && lex.data[p-2] != '\\' {
		return true
	}

	return !(lex.data[p] == s)
}

func (lex *Lexer) isHeredocEnd(p int) bool {
	if lex.data[p-1] != '\r' && lex.data[p-1] != '\n' {
		return false
	}

	if p == len(lex.data) {
		return false
	}

	for lex.data[p] == ' ' || lex.data[p] == '\t' {
		p++
	}

	l := len(lex.heredocLabel)
	if len(lex.data) < p+l {
		return false
	}

	if len(lex.data) > p+l && isValidVarName(lex.data[p+l]) {
		return false
	}

	a := string(lex.heredocLabel)
	b := string(lex.data[p : p+l])

	_, _ = a, b

	if bytes.Equal(lex.heredocLabel, lex.data[p:p+l]) {
		lex.p = p
		return true
	}

	return false
}

func (lex *Lexer) isNotHeredocEnd(p int) bool {
	return !lex.isHeredocEnd(p)
}

func (lex *Lexer) growCallStack() {
	if lex.top == len(lex.stack) {
		lex.stack = append(lex.stack, 0)
	}
}

func (lex *Lexer) isNotPhpCloseToken() bool {
	if lex.p+1 == len(lex.data) {
		return true
	}

	return lex.data[lex.p] != '?' || lex.data[lex.p+1] != '>'
}

func (lex *Lexer) isNotNewLine() bool {
	if lex.data[lex.p] == '\n' && lex.data[lex.p-1] == '\r' {
		return true
	}

	return lex.data[lex.p-1] != '\n' && lex.data[lex.p-1] != '\r'
}

func (lex *Lexer) call(state int, fnext int) {
	lex.growCallStack()

	lex.stack[lex.top] = state
	lex.top++

	lex.p++
	lex.cs = fnext
}

func (lex *Lexer) ret(n int) {
	lex.top = lex.top - n
	if lex.top < 0 {
		lex.top = 0
	}
	lex.cs = lex.stack[lex.top]
	lex.p++
}

func (lex *Lexer) ungetStr(s string) {
	tokenStr := string(lex.data[lex.ts:lex.te])
	if strings.HasSuffix(tokenStr, s) {
		lex.ungetCnt(len(s))
	}
}

func (lex *Lexer) ungetFromStart(n int) {
	tokenLength := lex.te - lex.ts
	lex.ungetCnt(tokenLength - n)
}

func (lex *Lexer) ungetCnt(n int) {
	lex.p = lex.p - n
	lex.te = lex.te - n
}

func (lex *Lexer) ungetWhile(s byte) {
	for i := 0; i < 100; i++ {
		v := lex.data[lex.te]
		if v == s {
			break
		}

		lex.te--
		lex.p--
	}

	lex.te++
	lex.p++
}

func (lex *Lexer) error(msg string) {
	if lex.errHandlerFunc == nil {
		return
	}

	sl, slb := lex.newLines.GetLine(lex.ts)
	el, elb := lex.newLines.GetLine(lex.te - 1)
	pos := position.NewPosition(
		sl,
		el,
		lex.ts,
		lex.te,
		lex.ts-slb,
		lex.te-elb,
	)

	lex.errHandlerFunc(errors.NewError(msg, pos))
}

func isValidVarNameStart(r byte) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' || r >= 0x80
}

func isValidVarName(r byte) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' ||
		r >= 0x80
}

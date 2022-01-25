package compiler

import (
	. "glox-vm"
	"reflect"
	"strconv"
)

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

type (
	Precedence uint8
	parseRule  struct {
		prefix string
		infix  string
		pcd    Precedence
	}
)

var rules = [40]parseRule{}

func init() {
	rules[TOKEN_LEFT_PAREN] = parseRule{"Grouping", "", PREC_NONE}
	rules[TOKEN_RIGHT_PAREN] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_LEFT_BRACE] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_RIGHT_BRACE] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_COMMA] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_DOT] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_MINUS] = parseRule{"Unary", "Binary", PREC_TERM}
	rules[TOKEN_PLUS] = parseRule{"", "Binary", PREC_TERM}
	rules[TOKEN_SEMICOLON] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_SLASH] = parseRule{"", "Binary", PREC_FACTOR}
	rules[TOKEN_STAR] = parseRule{"", "Binary", PREC_FACTOR}
	rules[TOKEN_BANG] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_BANG_EQUAL] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_EQUAL] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_EQUAL_EQUAL] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_GREATER] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_GREATER_EQUAL] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_LESS] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_LESS_EQUAL] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_IDENTIFIER] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_STRING] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_NUMBER] = parseRule{"Number", "", PREC_NONE}
	rules[TOKEN_AND] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_CLASS] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_ELSE] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_FALSE] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_FOR] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_FUN] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_IF] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_NIL] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_OR] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_PRINT] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_RETURN] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_SUPER] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_THIS] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_TRUE] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_VAR] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_WHILE] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_ERROR] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_EOF] = parseRule{"", "", PREC_NONE}
}

func getRule(tokenType uint8) *parseRule {
	return &rules[tokenType]
}

func (compiler *Compiler) callParseFn(parseFn string) {
	if parseFn == "" {
		return
	}
	value := reflect.ValueOf(compiler)
	method := value.MethodByName(parseFn)
	method.Call([]reflect.Value{})
}

// 解析相等或更高优先级的符号
func (compiler *Compiler) parsePrecedence(pcd Precedence) {
	compiler.advance()
	prefixFn := getRule(compiler.parser.previous.tokenType).prefix
	if prefixFn == "" {
		compiler.parser.errorAtPrevious("Expect expression.")
		return
	}
	compiler.callParseFn(prefixFn)
	for pcd <= getRule(compiler.parser.current.tokenType).pcd {
		compiler.advance()
		infixFn := getRule(compiler.parser.previous.tokenType).infix
		compiler.callParseFn(infixFn)
	}
}

func (compiler *Compiler) expression() {
	compiler.parsePrecedence(PREC_ASSIGNMENT)
}

func (compiler *Compiler) Number() {
	d, _ := strconv.ParseFloat(string(compiler.parser.previous.lexeme), 64)
	compiler.emitConstant(Value(d))
}

func (compiler *Compiler) Grouping() {
	compiler.expression()
	compiler.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (compiler *Compiler) Unary() {
	typ := compiler.parser.previous.tokenType
	compiler.parsePrecedence(PREC_UNARY)
	switch typ {
	case TOKEN_MINUS:
		compiler.emit(OP_NEGATE)
	default:
		return
	}
}

func (compiler *Compiler) Binary() {
	typ := compiler.parser.previous.tokenType
	rule := getRule(typ)
	compiler.parsePrecedence(rule.pcd + 1)
	switch typ {
	case TOKEN_PLUS:
		compiler.emit(OP_ADD)
	case TOKEN_MINUS:
		compiler.emit(OP_SUBTRACT)
	case TOKEN_STAR:
		compiler.emit(OP_MULTIPLY)
	case TOKEN_SLASH:
		compiler.emit(OP_DIVIDE)
	default:
		return
	}
}

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
	rules[TOKEN_BANG] = parseRule{"Unary", "", PREC_NONE}
	rules[TOKEN_BANG_EQUAL] = parseRule{"", "Binary", PREC_EQUALITY}
	rules[TOKEN_EQUAL] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_EQUAL_EQUAL] = parseRule{"", "Binary", PREC_EQUALITY}
	rules[TOKEN_GREATER] = parseRule{"", "Binary", PREC_COMPARISON}
	rules[TOKEN_GREATER_EQUAL] = parseRule{"", "Binary", PREC_COMPARISON}
	rules[TOKEN_LESS] = parseRule{"", "Binary", PREC_COMPARISON}
	rules[TOKEN_LESS_EQUAL] = parseRule{"", "Binary", PREC_COMPARISON}
	rules[TOKEN_IDENTIFIER] = parseRule{"Variable", "", PREC_NONE}
	rules[TOKEN_STRING] = parseRule{"String", "", PREC_NONE}
	rules[TOKEN_NUMBER] = parseRule{"Number", "", PREC_NONE}
	rules[TOKEN_AND] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_CLASS] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_ELSE] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_FALSE] = parseRule{"Literal", "", PREC_NONE}
	rules[TOKEN_FOR] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_FUN] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_IF] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_NIL] = parseRule{"Literal", "", PREC_NONE}
	rules[TOKEN_OR] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_PRINT] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_RETURN] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_SUPER] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_THIS] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_TRUE] = parseRule{"Literal", "", PREC_NONE}
	rules[TOKEN_VAR] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_WHILE] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_ERROR] = parseRule{"", "", PREC_NONE}
	rules[TOKEN_EOF] = parseRule{"", "", PREC_NONE}
}

func getRule(tokenType uint8) *parseRule {
	return &rules[tokenType]
}

func (compiler *Compiler) callParseFn(parseFn string, canAssign bool) {
	if parseFn == "" {
		return
	}
	value := reflect.ValueOf(compiler)
	method := value.MethodByName(parseFn)
	method.Call([]reflect.Value{reflect.ValueOf(canAssign)})
}

// 解析相等或更高优先级的符号
func (compiler *Compiler) parsePrecedence(pcd Precedence) {
	compiler.advance()
	prefixFn := getRule(compiler.parser.previous.tokenType).prefix
	if prefixFn == "" {
		compiler.parser.errorAtPrevious("Expect expression.")
		return
	}
	canAssign := pcd <= PREC_ASSIGNMENT
	compiler.callParseFn(prefixFn, canAssign)
	for pcd <= getRule(compiler.parser.current.tokenType).pcd {
		compiler.advance()
		infixFn := getRule(compiler.parser.previous.tokenType).infix
		compiler.callParseFn(infixFn, canAssign)
	}
	if !canAssign && compiler.match(TOKEN_EQUAL) {
		compiler.parser.errorAtPrevious("Invalid assignment target.")
	}
}

func (compiler *Compiler) expression() {
	compiler.parsePrecedence(PREC_ASSIGNMENT)
}

func (compiler *Compiler) Number(bool) {
	d, _ := strconv.ParseFloat(string(compiler.parser.previous.lexeme), 64)
	compiler.emitOpConstant(NewNumber(d))
}

func (compiler *Compiler) Grouping(bool) {
	compiler.expression()
	compiler.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (compiler *Compiler) Unary(bool) {
	typ := compiler.parser.previous.tokenType
	compiler.parsePrecedence(PREC_UNARY)
	switch typ {
	case TOKEN_MINUS:
		compiler.emit(OP_NEGATE)
	case TOKEN_BANG:
		compiler.emit(OP_NOT)
	default:
		return
	}
}

func (compiler *Compiler) Binary(bool) {
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
	case TOKEN_BANG_EQUAL:
		compiler.emit(OP_EQUAL, OP_NOT)
	case TOKEN_EQUAL_EQUAL:
		compiler.emit(OP_EQUAL)
	case TOKEN_GREATER:
		compiler.emit(OP_GREATER)
	case TOKEN_GREATER_EQUAL:
		compiler.emit(OP_LESS, OP_NOT)
	case TOKEN_LESS:
		compiler.emit(OP_LESS)
	case TOKEN_LESS_EQUAL:
		compiler.emit(OP_GREATER, OP_NOT)
	default:
		return
	}
}

func (compiler *Compiler) Literal(bool) {
	switch compiler.parser.previous.tokenType {
	case TOKEN_FALSE: compiler.emit(OP_FALSE)
	case TOKEN_TRUE: compiler.emit(OP_TRUE)
	case TOKEN_NIL: compiler.emit(OP_NIL)
	default:
		return
	}
}

func (compiler *Compiler) String(bool) {
	str := string(compiler.parser.previous.lexeme)
	compiler.emitOpConstant(NewString(str))
}

func (compiler *Compiler) Variable(canAssign bool) {
	compiler.namedVariable(compiler.parser.previous, canAssign)
}

func (compiler *Compiler) namedVariable(name Token, canAssign bool) {
	arg := compiler.identifierConstant(&name)
	if canAssign && compiler.match(TOKEN_EQUAL) {
		compiler.expression()
		compiler.emit(OP_SET_GLOBAL, arg)
	} else {
		compiler.emit(OP_GET_GLOBAL, arg)
	}
}
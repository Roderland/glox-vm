package compiler

import (
	"github.com/Roderland/glox-vm/chunk"
	"strconv"
)

type Precedence uint8

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

type parseRule struct {
	prefix func()
	infix  func()
	pd     Precedence
}

var rules = [40]parseRule{}

func init() {
	rules[TOKEN_LEFT_PAREN] = parseRule{grouping, call, PREC_CALL}
	rules[TOKEN_RIGHT_PAREN] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_LEFT_BRACE] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_RIGHT_BRACE] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_COMMA] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_DOT] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_MINUS] = parseRule{unary, binary, PREC_TERM}
	rules[TOKEN_PLUS] = parseRule{nil, binary, PREC_TERM}
	rules[TOKEN_SEMICOLON] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SLASH] = parseRule{nil, binary, PREC_FACTOR}
	rules[TOKEN_STAR] = parseRule{nil, binary, PREC_FACTOR}
	rules[TOKEN_BANG] = parseRule{unary, nil, PREC_NONE}
	rules[TOKEN_BANG_EQUAL] = parseRule{nil, binary, PREC_EQUALITY}
	rules[TOKEN_EQUAL] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_EQUAL_EQUAL] = parseRule{nil, binary, PREC_EQUALITY}
	rules[TOKEN_GREATER] = parseRule{nil, binary, PREC_COMPARISON}
	rules[TOKEN_GREATER_EQUAL] = parseRule{nil, binary, PREC_COMPARISON}
	rules[TOKEN_LESS] = parseRule{nil, binary, PREC_COMPARISON}
	rules[TOKEN_LESS_EQUAL] = parseRule{nil, binary, PREC_COMPARISON}
	rules[TOKEN_IDENTIFIER] = parseRule{variable, nil, PREC_NONE}
	rules[TOKEN_STRING] = parseRule{str, nil, PREC_NONE}
	rules[TOKEN_NUMBER] = parseRule{number, nil, PREC_NONE}
	rules[TOKEN_AND] = parseRule{nil, and, PREC_AND}
	rules[TOKEN_CLASS] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ELSE] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FALSE] = parseRule{literal, nil, PREC_NONE}
	rules[TOKEN_FOR] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FUN] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_IF] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_NIL] = parseRule{literal, nil, PREC_NONE}
	rules[TOKEN_OR] = parseRule{nil, or, PREC_OR}
	rules[TOKEN_PRINT] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_RETURN] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SUPER] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_THIS] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_TRUE] = parseRule{literal, nil, PREC_NONE}
	rules[TOKEN_VAR] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_WHILE] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ERROR] = parseRule{nil, nil, PREC_NONE}
	rules[TOKEN_EOF] = parseRule{nil, nil, PREC_NONE}
}

func or()   {}
func and()  {}
func str()  {}
func call() {}
func literal() {
	tp := prs.previous.tp
	switch tp {
	case TOKEN_NIL:
		emitBytes(chunk.OP_NIL)
	case TOKEN_FALSE:
		emitBytes(chunk.OP_FALSE)
	case TOKEN_TRUE:
		emitBytes(chunk.OP_TRUE)
	default:
		return
	}
}
func variable() {}

func expression() {
	parsePrecedence(PREC_ASSIGNMENT)
}

func parsePrecedence(pd Precedence) {
	advance()
	prefixFn := getParseRule(prs.previous.tp).prefix
	if prefixFn == nil {
		errorAtPrevious("Expect expression.")
		return
	}
	prefixFn()

	for pd <= getParseRule(prs.current.tp).pd {
		advance()
		infixFn := getParseRule(prs.previous.tp).infix
		infixFn()
	}
}

func number() {
	float, _ := strconv.ParseFloat(prs.previous.lexeme, 64)
	emitConstant(chunk.NewNumber(float))
}

func grouping() {
	expression()
	consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func unary() {
	operatorType := prs.previous.tp

	parsePrecedence(PREC_UNARY)

	switch operatorType {
	case TOKEN_MINUS:
		emitBytes(chunk.OP_NEGATE)
	case TOKEN_BANG:
		emitBytes(chunk.OP_NOT)
	default:
		return
	}
}

func getParseRule(tp tokenType) *parseRule {
	return &rules[tp]
}

func binary() {
	operatorType := prs.previous.tp

	parsePrecedence(getParseRule(operatorType).pd + 1)

	switch operatorType {
	case TOKEN_PLUS:
		emitBytes(chunk.OP_ADD)
	case TOKEN_MINUS:
		emitBytes(chunk.OP_SUBTRACT)
	case TOKEN_STAR:
		emitBytes(chunk.OP_MULTIPLY)
	case TOKEN_SLASH:
		emitBytes(chunk.OP_DIVIDE)
	case TOKEN_BANG_EQUAL:
		emitBytes(chunk.OP_EQUAL, chunk.OP_NOT)
	case TOKEN_EQUAL_EQUAL:
		emitBytes(chunk.OP_EQUAL)
	case TOKEN_GREATER:
		emitBytes(chunk.OP_GREATER)
	case TOKEN_GREATER_EQUAL:
		emitBytes(chunk.OP_LESS, chunk.OP_NOT)
	case TOKEN_LESS:
		emitBytes(chunk.OP_LESS)
	case TOKEN_LESS_EQUAL:
		emitBytes(chunk.OP_GREATER, chunk.OP_NOT)
	default:
		return
	}
}

package compiler

import (
	. "glox-vm"
	"math"
)

type Compiler struct {
	scanner *Scanner
	parser  *Parser
	//chunk   *Chunk
	scope   *Scope
}

func (compiler *Compiler) Compile(source []byte) (*FuncData, error) {
	compiler.scanner = &Scanner{source: source, line: 1}
	compiler.parser = &Parser{}
	//compiler.chunk = CreateChunk()
	compiler.scope = &Scope{function: NewFuncData(TYPE_SCRIPT, ""), compiler: compiler}
	compiler.advance()
	for !compiler.match(TOKEN_EOF) {
		compiler.declaration()
	}
	return compiler.endCompile(), compiler.parser.err
}

func (compiler *Compiler) endCompile() *FuncData {
	function := compiler.scope.function
	if function.Ft == TYPE_SCRIPT {
		compiler.emit(OP_RETURN)
	}
	if compiler.parser.err == nil {
		funcName := function.Name
		if funcName == "" {
			funcName = "<script>"
		}
		DisassembleChunk(compiler.currentChunk(), funcName)
	}
	compiler.scope = compiler.scope.enclosing
	return function
}

func (compiler *Compiler) currentChunk() *Chunk {
	return &compiler.scope.function.FunChunk
}

func (compiler *Compiler) advance() {
	compiler.parser.previous = compiler.parser.current
	for {
		compiler.parser.current = compiler.scanner.scanToken()
		if compiler.parser.current.tokenType != TOKEN_ERROR {
			break
		}
		compiler.parser.errorAtCurrent(string(compiler.parser.current.lexeme))
	}
}

func (compiler *Compiler) match(typ uint8) bool {
	if !compiler.check(typ) {
		return false
	}
	compiler.advance()
	return true
}

func (compiler *Compiler) check(typ uint8) bool {
	return compiler.parser.current.tokenType == typ
}

func (compiler *Compiler) consume(tokenType uint8, msg string) {
	if compiler.parser.current.tokenType == tokenType {
		compiler.advance()
		return
	}
	compiler.parser.errorAtCurrent(msg)
}

func (compiler *Compiler) synchronize() {
	compiler.parser.panicMode = false
	for compiler.parser.current.tokenType != TOKEN_EOF {
		if compiler.parser.current.tokenType == TOKEN_SEMICOLON {
			return
		}
		switch compiler.parser.current.tokenType {
		case TOKEN_CLASS:
			return
		case TOKEN_FUN:
			return
		case TOKEN_VAR:
			return
		case TOKEN_FOR:
			return
		case TOKEN_IF:
			return
		case TOKEN_WHILE:
			return
		case TOKEN_PRINT:
			return
		case TOKEN_RETURN:
			return
		default:

		}
		compiler.advance()
	}
}

func (compiler *Compiler) declaration() {
	if compiler.match(TOKEN_VAR) {
		compiler.varDeclaration()
	} else if compiler.match(TOKEN_FUN) {
		compiler.funDeclaration()
	} else {
		compiler.statement()
	}
}

func (compiler *Compiler) statement() {
	if compiler.match(TOKEN_PRINT) {
		compiler.printStatement()
	} else if compiler.match(TOKEN_IF) {
		compiler.ifStatement()
	} else if compiler.match(TOKEN_WHILE) {
		compiler.whileStatement()
	} else if compiler.match(TOKEN_LEFT_BRACE) {
		compiler.scope.begin()
		compiler.block()
		compiler.scope.end()
	} else if compiler.match(TOKEN_FOR) {
		compiler.forStatement()
	} else if compiler.match(TOKEN_RETURN) {
		compiler.returnStatement()
	} else {
		compiler.expressionStatement()
	}
}

// ifStatement if跳转的字节码实现：
// 这三个byte用于跳过OP_POP,then语句,OP_JUMP		OP_JUMP_IF_FALSE 	high-8bit low-8bit
// 将if条件值从VM栈中弹出							OP_POP
//												... then ...
// 这三个byte用于跳过OP_POP,else语句				OP_JUMP				high-8bit low-8bit
// 将if条件值从VM栈中弹出							OP_POP
// 												... else ...
func (compiler *Compiler) ifStatement() {
	compiler.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'if'.")
	compiler.expression()
	compiler.consume(TOKEN_RIGHT_PAREN, "Expect ')' after condition.")
	// OP_JUMP_IF_FALSE
	thenJump := compiler.emitJump(OP_JUMP_IF_FALSE)
	// 条件为真，将if条件值从VM栈中弹出
	compiler.emit(OP_POP)
	// ... then ...
	compiler.statement()
	// OP_JUMP
	elseJump := compiler.emitJump(OP_JUMP)
	// 补充OP_JUMP_IF_FALSE的16位跳转地址
	compiler.patchJump(thenJump)
	// 条件为假，将if条件值从VM栈中弹出
	compiler.emit(OP_POP)
	// ... else ...
	if compiler.match(TOKEN_ELSE) {
		compiler.statement()
	}
	// 补充OP_JUMP的16位跳转地址
	compiler.patchJump(elseJump)
}

// whileStatement while循环的字节码实现：
// 												... condition ...
// condition为false跳出循环						OP_JUMP_IF_FALSE	high-8bit low-8bit
// 												OP_POP
// 												... body ...
// body执行完跳转至condition						OP_LOOP				high-8bit low-8bit
// 												OP_POP
func (compiler *Compiler) whileStatement() {
	// 记录循环起点用于跳回
	loopStart := len(compiler.currentChunk().Bytecodes)
	compiler.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'while'.")
	compiler.expression()
	compiler.consume(TOKEN_RIGHT_PAREN, "Expect ')' after condition.")
	// OP_JUMP_IF_FALSE
	exitJump := compiler.emitJump(OP_JUMP_IF_FALSE)
	// 条件为真，将while条件值从VM栈中弹出
	compiler.emit(OP_POP)
	// while循环体
	compiler.statement()
	// OP_LOOP
	compiler.emitLoop(loopStart)
	// 补充OP_JUMP_IF_FALSE的16位跳转地址
	compiler.patchJump(exitJump)
	// 条件为假，将while条件值从VM栈中弹出
	compiler.emit(OP_POP)
}

// forStatement for循环的字节码实现：
// 	==> condition执行完为false跳出循环
// 	==> condition执行完为true跳转至body
//	==> body执行完跳转至increment
//	==> increment执行完跳转至condition
// 												... initializer ...
// 												... condition ...
// 	condition执行完为false跳出循环				OP_JUMP_IF_FALSE	high-8bit low-8bit
// 												OP_POP
// 	condition执行完为true跳转至body				OP_JUMP				high-8bit low-8bit
// 												... increment ...
// 												OP_POP
// 	increment执行完跳转至condition				OP_LOOP				high-8bit low-8bit
// 												... body ...
// 	body执行完跳转至increment					OP_LOOP				high-8bit low-8bit
// 												OP_POP
func (compiler *Compiler) forStatement() {
	compiler.scope.begin()
	compiler.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'for'.")
	// ... initializer ...
	if compiler.match(TOKEN_SEMICOLON) {
		// No initializer.
	} else if compiler.match(TOKEN_VAR) {
		compiler.varDeclaration()
	} else {
		compiler.expressionStatement()
	}
	exitJump := -1
	conditionStart := len(compiler.currentChunk().Bytecodes)
	incrementStart := conditionStart
	if !compiler.match(TOKEN_SEMICOLON) {
		// ... condition ...
		compiler.expression()
		compiler.consume(TOKEN_SEMICOLON, "Expect ';' after loop condition.")
		// OP_JUMP_IF_FALSE用于退出循环
		exitJump = compiler.emitJump(OP_JUMP_IF_FALSE)
		// 条件为真,将for条件值从VM栈中弹出
		compiler.emit(OP_POP)
	}
	if !compiler.match(TOKEN_RIGHT_PAREN) {
		// OP_JUMP用于跳转到body先执行
		bodyJump := compiler.emitJump(OP_JUMP)
		// 存在increment表达式时使body执行完跳转至increment
		incrementStart = len(compiler.currentChunk().Bytecodes)
		// ... increment ...
		compiler.expression()
		// 将increment结果值从VM栈中弹出
		compiler.emit(OP_POP)
		compiler.consume(TOKEN_RIGHT_PAREN, "Expect ')' after for clauses.")
		// OP_LOOP用于跳转到循环... condition ...判断
		compiler.emitLoop(conditionStart)
		// 补充OP_JUMP跳转到body的16位跳转地址
		compiler.patchJump(bodyJump)
	}
	// ... body ...
	compiler.statement()
	// 循环体结束后跳转至 ... increment ...
	compiler.emitLoop(incrementStart)
	if exitJump != -1 {
		// 补充OP_JUMP_IF_FALSE退出for循环的16位跳转地址
		compiler.patchJump(exitJump)
		// 条件为假,将for条件值从VM栈中弹出
		compiler.emit(OP_POP)
	}
	compiler.scope.end()
}

// patchJump 补充jump操作的16位跳转地址
func (compiler *Compiler) patchJump(offset int) {
	jump := len(compiler.currentChunk().Bytecodes) - offset - 2
	if jump > math.MaxUint8 {
		compiler.parser.errorAtPrevious("Too much code to jump over.")
	}
	compiler.currentChunk().Bytecodes[offset] = uint8(jump>>8) & 0xff
	compiler.currentChunk().Bytecodes[offset+1] = uint8(jump) & 0xff
}

func (compiler *Compiler) block() {
	for !compiler.check(TOKEN_RIGHT_BRACE) && !compiler.check(TOKEN_EOF) {
		compiler.declaration()
	}
	compiler.consume(TOKEN_RIGHT_BRACE, "Expect '}' after block.")
}

func (compiler *Compiler) printStatement() {
	compiler.expression()
	compiler.consume(TOKEN_SEMICOLON, "Expect ';' after value.")
	compiler.emit(OP_PRINT)
}

func (compiler *Compiler) expressionStatement() {
	compiler.expression()
	compiler.consume(TOKEN_SEMICOLON, "Expect ';' after expression.")
	compiler.emit(OP_POP)
}

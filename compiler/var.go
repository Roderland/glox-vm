package compiler

import (
	. "glox-vm"
)

/*
	以下是compiler对变量声明和定义的过程
	注意：全局变量在常量池中添加变量记录; 局部变量在作用域的局部变量表中添加记录。
*/
func (compiler *Compiler) varDeclaration() {
	variable := compiler.parseVariable("Expect variable name.")
	if compiler.match(TOKEN_EQUAL) {
		compiler.expression()
	} else {
		compiler.emit(OP_NIL)
	}
	compiler.consume(TOKEN_SEMICOLON, "Expect ';' after variable declaration.")
	compiler.defineVariable(variable)
}

// parseVariable 进行变量的解析：
// 1.声明变量;
// 2.如果是全局变量则向常量池添加变量名字符串并返回索引;
func (compiler *Compiler) parseVariable(msg string) uint8 {
	compiler.consume(TOKEN_IDENTIFIER, msg)
	compiler.declareVariable()
	if compiler.scope.scopeDepth > 0 {
		return 0
	}
	return compiler.identifierConstant(&compiler.parser.previous)
}

// declareVariable 进行变量声明：
// 如果是全局变量直接返回;
// 如果是局部变量则在作用域的局部变量表中添加变量的记录;
func (compiler *Compiler) declareVariable() {
	if compiler.scope.scopeDepth == 0 {
		return
	}
	name := &compiler.parser.previous
	for i := compiler.scope.localCount - 1; i >= 0; i-- {
		local := &compiler.scope.locals[i]
		// 如果local的depth已经小于当前scope的depth，已经说明变量在当前作用域下未经过声明
		if local.depth != -1 && local.depth < compiler.scope.scopeDepth {
			break
		}
		// 当前作用域下已经对该变量进行过声明，给出error
		if bytesEqual(name.lexeme, local.name.lexeme) {
			compiler.parser.errorAtPrevious("Already a variable with this name in this scope.")
		}
	}
	// 在作用域的局部变量表中添加变量的记录
	compiler.scope.addLocal(*name)
}

// identifierConstant 向常量池添加变量名字符串并返回索引
func (compiler *Compiler) identifierConstant(name *Token) uint8 {
	return compiler.emitConstant(NewString(string(name.lexeme)))
}

// defineVariable 进行变量定义：
// 对于全局变量：向chunk输出全局变量的定义操作字节码和全局变量名称在常量池的索引
// 对于局部变量：标记为已初始化
func (compiler *Compiler) defineVariable(global uint8) {
	if compiler.scope.scopeDepth > 0 {
		compiler.scope.markInitialized()
		return
	}
	compiler.emit(OP_DEFINE_GLOBAL, global)
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := range a {
		if b[i] != c {
			return false
		}
	}
	return true
}


package compiler

import . "glox-vm"

/*
	以下是compiler对函数声明和定义的过程
	注意：
	在前面的章节中，我解释了变量是如何分两个阶段定义的。这确保您无法在变量自己的初始化程序中访问变量的值。那会很糟糕，因为变量还没有值。
	函数不会遇到这个问题。函数在其主体中引用自己的名称是安全的。
	在完全定义之前，您不能调用函数并执行主体，因此您永远不会看到处于未初始化状态的变量。
	实际上，允许这样做以支持递归本地函数很有用。
	为了使它工作，我们在编译名称之前，在编译主体之前，将函数声明的变量标记为“已初始化”。这样，名称就可以在正文中被引用而不会产生错误。
*/
func (compiler *Compiler) funDeclaration() {
	variable := compiler.parseVariable("Expect function name.")
	compiler.scope.markInitialized()
	compiler.function(TYPE_FUNCTION)
	compiler.defineVariable(variable)
}

func (compiler *Compiler) function(ft FuncType) {
	compiler.declareVariable()
	name := string(compiler.parser.previous.lexeme)
	scope := &Scope{function: NewFuncData(ft, name), enclosing: compiler.scope, compiler: compiler}
	scope.compiler.scope = scope
	scope.begin()
	compiler.consume(TOKEN_LEFT_PAREN, "Expect '(' after function name.")
	if !compiler.check(TOKEN_RIGHT_PAREN) {
		for {
			scope.function.Arity ++
			if scope.function.Arity > 255 {
				compiler.parser.errorAtCurrent("Can't have more than 255 parameters.");
			}
			param := compiler.parseVariable("Expect parameter name.")
			compiler.defineVariable(param)
			if !compiler.match(TOKEN_COMMA) {
				break
			}
		}
	}
	compiler.consume(TOKEN_RIGHT_PAREN, "Expect ')' after parameters.")
	compiler.consume(TOKEN_LEFT_BRACE, "Expect '{' before function body.")
	compiler.block()
	funData := compiler.endCompile()
	compiler.emitOpConstant(NewFunction(*funData))
}

func (compiler *Compiler) returnStatement() {
	if compiler.scope.function.Ft == TYPE_SCRIPT {
		compiler.parser.errorAtPrevious("Can't return from top-level code.")
	}
	if compiler.match(TOKEN_SEMICOLON) {
		compiler.emit(OP_NIL, OP_RETURN)
	} else {
		compiler.expression()
		compiler.consume(TOKEN_SEMICOLON, "Expect ';' after return value.")
		compiler.emit(OP_RETURN)
	}
}
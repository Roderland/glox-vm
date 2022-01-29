package compiler

import (
	. "glox-vm"
	"math"
)

const MAX_LOCAL_COUNT = math.MaxUint8 + 1

type (
	Local struct {
		name  Token
		depth int
	}
	Scope struct {
		locals     [MAX_LOCAL_COUNT]Local // 局部变量表
		localCount int                    // 局部变量表大小
		scopeDepth int
		compiler   *Compiler
	}
)

// begin 进入一个作用域
func (scope *Scope) begin() {
	scope.scopeDepth++
}

// end 退出一个作用域
// 将被退出的作用域包含的VM栈帧退出
func (scope *Scope) end() {
	scope.scopeDepth--
	for scope.localCount > 0 && scope.locals[scope.localCount-1].depth > scope.scopeDepth {
		scope.compiler.emit(OP_POP)
		scope.localCount--
	}
}

// addLocal 在作用域的局部变量表中添加变量记录
func (scope *Scope) addLocal(name Token) {
	if scope.localCount == MAX_LOCAL_COUNT {
		scope.compiler.parser.errorAtPrevious("Too many local variables in function.")
		return
	}
	scope.locals[scope.localCount] = Local{
		name:  name,
		depth: -1,
	}
	scope.localCount++
}

// resolveLocal 解析局部变量：
// 返回局部变量在局部变量表中的索引
// 对于不存在的局部变量返回 false
func (scope *Scope) resolveLocal(name *Token) (uint8, bool) {
	for i := scope.localCount - 1; i >= 0; i-- {
		local := scope.locals[i]
		if bytesEqual(name.lexeme, local.name.lexeme) {
			// 对于未初始化的局部变量的解析将发生 error
			if local.depth == -1 {
				scope.compiler.parser.errorAtPrevious("Can't read local variable in its own initializer.")
			}
			return uint8(i), true
		}
	}
	return 0, false
}

// markInitialized 初始化标记：
// 将局部变量表末尾的变量标记为已初始化
func (scope *Scope) markInitialized() {
	scope.locals[scope.localCount-1].depth = scope.scopeDepth
}

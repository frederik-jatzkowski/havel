package parser

type Visitor interface {
	// program
	VisitProgram(program *Program)

	// packages
	VisitPackage(pkg *Package)

	// functions
	VisitFunction(function *Function)
	VisitFunctionVariableDeclaration(declaration *FunctionVariableDeclaration)

	// basic blocks
	VisitBlock(block *BasicBlock)
	VisitJump(terminator *Jump)
	VisitConditionalJump(terminator *ConditionalJump)
	VisitReturn(terminator *Return)

	// instructions
	VisitInstruction(instr *Instruction)
	VisitWriteRegister(write *WriteRegister)
	VisitWriteVariable(write *WriteVariable)
	VisitReadRegister(read *ReadRegister)
	VisitReadVariable(read *ReadVariable)

	// operations
	VisitAluOperation(op *AluOperation)
	VisitLocalCall(op *LocalCall)
	VisitDebugOperation(op *DebugOperation)
	VisitPrimitiveLiteral(op *PrimitiveLiteral)

	// visiting context
	SetCurrentProgram(program *Program)
	SetCurrentPackage(pkg *Package)
	SetCurrentFunction(function *Function)
	SetCurrentBlock(block *BasicBlock)
}

type VisitorContext struct {
	CurrentProgram  *Program
	CurrentPackage  *Package
	CurrentFunction *Function
	CurrentBlock    *BasicBlock
}

func (ctx *VisitorContext) SetCurrentProgram(program *Program) {
	ctx.CurrentProgram = program
}

func (ctx *VisitorContext) SetCurrentPackage(pkg *Package) {
	ctx.CurrentPackage = pkg
}

func (ctx *VisitorContext) SetCurrentFunction(function *Function) {
	ctx.CurrentFunction = function
}

func (ctx *VisitorContext) SetCurrentBlock(block *BasicBlock) {
	ctx.CurrentBlock = block
}

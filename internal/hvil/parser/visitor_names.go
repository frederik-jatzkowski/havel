package parser

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type NameResolution struct {
	VisitorContext
	Result *errors.Collector
}

func (pass *NameResolution) VisitProgram(program *Program) {
	program.packageMap = make(map[string]*Package, len(program.Packages))
	for _, pkg := range program.Packages {
		_, exists := program.packageMap[pkg.Name]
		if exists {
			pass.Result.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is redeclared in this program", pkg.Name),
			)
		}

		program.packageMap[pkg.Name] = pkg
	}
}

func (pass *NameResolution) VisitPackage(pkg *Package) {
	pkg.functionsMap = make(map[string]*Function, len(pkg.Functions))
	for _, function := range pkg.Functions {
		_, exists := pkg.functionsMap[function.Name]
		if exists {
			pass.Result.Err(
				function.Pos,
				"NameError",
				fmt.Sprintf("the function %s is redeclared in this package", function.Name),
			)
		}

		pkg.functionsMap[function.Name] = function
	}

	if pkg.IsMain {
		_, mainExists := pkg.functionsMap["main"]
		if !mainExists {
			pass.Result.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is the main package and does not contain a main function", pkg.Name),
			)
		}
	} else {
		_, mainExists := pkg.functionsMap["main"]
		if mainExists {
			pass.Result.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is not the main package but contains a main function", pkg.Name),
			)
		}
	}
}

func (pass *NameResolution) VisitFunction(function *Function) {
	mapSize := len(function.Parameters.Items) + len(function.LocalDeclarations.Items)
	if function.ReturnValue != nil {
		mapSize++
	}

	function.variableDeclarationMap = make(map[string]*FunctionVariableDeclaration, mapSize)

	function.blockMap = make(map[string]*BasicBlock, len(function.BasicBlocks))
	for _, block := range function.BasicBlocks {
		_, exists := function.blockMap[block.Identifier]
		if exists {
			pass.Result.Err(
				block.Pos,
				"NameError",
				fmt.Sprintf("the basic block %s is redeclared in this function", block.Identifier),
			)
		}

		function.blockMap[block.Identifier] = block
	}
}

func (pass *NameResolution) VisitFunctionVariableDeclaration(declaration *FunctionVariableDeclaration) {
	_, exists := declaration.function.variableDeclarationMap[declaration.Name]
	if exists {
		pass.Result.Err(
			declaration.Pos,
			"NameError",
			fmt.Sprintf("the variable %s is redeclared in this function", declaration.Name),
		)
	}

	declaration.function.variableDeclarationMap[declaration.Name] = declaration
}

func (pass *NameResolution) VisitBlock(block *BasicBlock) {
	block.registerMap = make(map[string]*WriteRegister, len(block.Instructions))
}
func (pass *NameResolution) VisitConditionalJump(terminator *ConditionalJump) {
	_, exists := terminator.block.function.blockMap[terminator.True]
	if !exists {
		pass.Result.Err(
			terminator.Pos,
			"NameError",
			fmt.Sprintf("the block %s does not exist in this function", terminator.True),
		)
	}

	_, exists = terminator.block.function.blockMap[terminator.False]
	if !exists {
		pass.Result.Err(
			terminator.Pos,
			"NameError",
			fmt.Sprintf("the block %s does not exist in this function", terminator.False),
		)
	}
}

func (pass *NameResolution) VisitJump(terminator *Jump) {
	_, exists := terminator.block.function.blockMap[terminator.Target]
	if !exists {
		pass.Result.Err(
			terminator.Pos,
			"NameError",
			fmt.Sprintf("the block %s does not exist in this function", terminator.Target),
		)
	}
}

func (pass *NameResolution) VisitReturn(terminator *Return) {}

func (pass *NameResolution) VisitInstruction(instr *Instruction) {}

func (pass *NameResolution) VisitWriteRegister(write *WriteRegister) {
	_, exists := write.block.registerMap[write.Identifier]
	if exists {
		pass.Result.Err(
			write.Pos,
			"NameError",
			fmt.Sprintf("the register %s is written twice in this basic block, this violates single static assignment form", write.Identifier),
		)
	}

	write.block.registerMap[write.Identifier] = write
}

func (pass *NameResolution) VisitWriteVariable(write *WriteVariable) {
	localDeclaration, exists := write.block.function.variableDeclarationMap[write.Identifier]
	if !exists {
		pass.Result.Err(
			write.Pos,
			"NameError",
			fmt.Sprintf("variable %s is not found in the current scope", write.Identifier),
		)
	} else {
		write.Declaration = localDeclaration
	}
}

func (pass *NameResolution) VisitReadRegister(read *ReadRegister) {
	declaration, exists := read.block.registerMap[read.Identifier]
	if !exists {
		pass.Result.Err(
			read.Pos,
			"NameError",
			fmt.Sprintf("the register %s is not yet defined in the current scope", read.Identifier),
		)
	} else {
		read.Declaration = declaration
	}
}

func (pass *NameResolution) VisitReadVariable(read *ReadVariable) {
	localDeclaration, exists := read.block.function.variableDeclarationMap[read.Identifier]
	if !exists {
		pass.Result.Err(
			read.Pos,
			"NameError",
			fmt.Sprintf("variable %s is not found in the current scope", read.Identifier),
		)
	} else {
		read.Declaration = localDeclaration
	}
}

func (pass *NameResolution) VisitAluOperation(op *AluOperation) {
	switch op.Name {
	case "add_u_32", "sub_u_32", "lt_u_32":
	default:
		pass.Result.Err(
			op.Pos,
			"NameError",
			fmt.Sprintf("there is no ALU operation called '%s'", op.Name),
		)
	}
}

func (pass *NameResolution) VisitLocalCall(op *LocalCall) {
	declaration, exists := op.block.function.pkg.functionsMap[op.Name]
	if !exists {
		pass.Result.Err(
			op.Pos,
			"NameError",
			fmt.Sprintf("the function %s is not defined locally in this package", op.Name),
		)
	} else {
		op.declaration = declaration
	}
}

func (pass *NameResolution) VisitDebugOperation(op *DebugOperation) {
	switch op.Name {
	case "print_u_32":
	default:
		pass.Result.Err(
			op.Pos,
			"NameError",
			fmt.Sprintf("there is no debug operation called '%s'", op.Name),
		)
	}
}

func (pass *NameResolution) VisitPrimitiveLiteral(op *PrimitiveLiteral) {}

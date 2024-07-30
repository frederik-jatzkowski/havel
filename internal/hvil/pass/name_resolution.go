package pass

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/internal/hvil/parser"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type NameResolution struct {
	Result *errors.Collector
	parser.VisitorContext
}

func (pass *NameResolution) VisitProgram(program *parser.Program) {
	program.PackageMap = make(map[string]*parser.Package, len(program.Packages))
	for _, pkg := range program.Packages {
		_, exists := program.PackageMap[pkg.Name]
		if exists {
			pass.Result.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is redeclared in this program", pkg.Name),
			)
		}

		program.PackageMap[pkg.Name] = pkg
	}
}

func (pass *NameResolution) VisitPackage(pkg *parser.Package) {
	pkg.FunctionsMap = make(map[string]*parser.Function, len(pkg.Functions))
	for _, function := range pkg.Functions {
		_, exists := pkg.FunctionsMap[function.Name]
		if exists {
			pass.Result.Err(
				function.Pos,
				"NameError",
				fmt.Sprintf("the function %s is redeclared in this package", function.Name),
			)
		}

		pkg.FunctionsMap[function.Name] = function
	}

	if pkg.IsMain {
		_, mainExists := pkg.FunctionsMap["main"]
		if !mainExists {
			pass.Result.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is the main package and does not contain a main function", pkg.Name),
			)
		}
	} else {
		_, mainExists := pkg.FunctionsMap["main"]
		if mainExists {
			pass.Result.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is not the main package but contains a main function", pkg.Name),
			)
		}
	}
}

func (pass *NameResolution) VisitFunction(function *parser.Function) {
	mapSize := len(function.Parameters.Items) + len(function.LocalDeclarations.Items)
	if function.ReturnValue != nil {
		mapSize++
	}

	function.VariableDeclarationMap = make(map[string]*parser.FunctionVariableDeclaration, mapSize)

	function.BlockMap = make(map[string]*parser.BasicBlock, len(function.BasicBlocks))
	for _, block := range function.BasicBlocks {
		_, exists := function.BlockMap[block.Identifier]
		if exists {
			pass.Result.Err(
				block.Pos,
				"NameError",
				fmt.Sprintf("the basic block %s is redeclared in this function", block.Identifier),
			)
		}

		function.BlockMap[block.Identifier] = block
	}
}

func (pass *NameResolution) VisitFunctionVariableDeclaration(declaration *parser.FunctionVariableDeclaration) {
	_, exists := pass.CurrentFunction.VariableDeclarationMap[declaration.Name]
	if exists {
		pass.Result.Err(
			declaration.Pos,
			"NameError",
			fmt.Sprintf("the variable %s is redeclared in this function", declaration.Name),
		)
	}

	pass.CurrentFunction.VariableDeclarationMap[declaration.Name] = declaration
}

func (pass *NameResolution) VisitBlock(block *parser.BasicBlock) {
	block.RegisterMap = make(map[string]*parser.WriteRegister, len(block.Instructions))
}
func (pass *NameResolution) VisitConditionalJump(terminator *parser.ConditionalJump) {
	_, exists := pass.CurrentFunction.BlockMap[terminator.True]
	if !exists {
		pass.Result.Err(
			terminator.Pos,
			"NameError",
			fmt.Sprintf("the block %s does not exist in this function", terminator.True),
		)
	}

	_, exists = pass.CurrentFunction.BlockMap[terminator.False]
	if !exists {
		pass.Result.Err(
			terminator.Pos,
			"NameError",
			fmt.Sprintf("the block %s does not exist in this function", terminator.False),
		)
	}
}

func (pass *NameResolution) VisitJump(terminator *parser.Jump) {
	_, exists := pass.CurrentFunction.BlockMap[terminator.Target]
	if !exists {
		pass.Result.Err(
			terminator.Pos,
			"NameError",
			fmt.Sprintf("the block %s does not exist in this function", terminator.Target),
		)
	}
}

func (pass *NameResolution) VisitReturn(terminator *parser.Return) {}

func (pass *NameResolution) VisitInstruction(instr *parser.Instruction) {}

func (pass *NameResolution) VisitWriteRegister(write *parser.WriteRegister) {
	_, exists := pass.CurrentBlock.RegisterMap[write.Identifier]
	if exists {
		pass.Result.Err(
			write.Pos,
			"NameError",
			fmt.Sprintf("the register %s is written twice in this basic block, this violates single static assignment form", write.Identifier),
		)
	}

	pass.CurrentBlock.RegisterMap[write.Identifier] = write
}

func (pass *NameResolution) VisitWriteVariable(write *parser.WriteVariable) {
	localDeclaration, exists := pass.CurrentFunction.VariableDeclarationMap[write.Identifier]
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

func (pass *NameResolution) VisitReadRegister(read *parser.ReadRegister) {
	declaration, exists := pass.CurrentBlock.RegisterMap[read.Identifier]
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

func (pass *NameResolution) VisitReadVariable(read *parser.ReadVariable) {
	localDeclaration, exists := pass.CurrentFunction.VariableDeclarationMap[read.Identifier]
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

func (pass *NameResolution) VisitAluOperation(op *parser.AluOperation) {
	_, exists := aluDefinitions[op.Name]
	if !exists {
		pass.Result.Err(
			op.Pos,
			"NameError",
			fmt.Sprintf("there is no alu operation called '%s'", op.Name),
		)
	}
}

func (pass *NameResolution) VisitLocalCall(op *parser.LocalCall) {
	declaration, exists := pass.CurrentPackage.FunctionsMap[op.Name]
	if !exists {
		pass.Result.Err(
			op.Pos,
			"NameError",
			fmt.Sprintf("the function %s is not defined locally in this package", op.Name),
		)
	} else {
		op.Declaration = declaration
	}
}

func (pass *NameResolution) VisitDebugOperation(op *parser.DebugOperation) {
	_, exists := debugDefinitions[op.Name]
	if !exists {
		pass.Result.Err(
			op.Pos,
			"NameError",
			fmt.Sprintf("there is no debug operation called '%s'", op.Name),
		)
	}
}

func (pass *NameResolution) VisitScalarLiteral(op *parser.ScalarLiteral) {}

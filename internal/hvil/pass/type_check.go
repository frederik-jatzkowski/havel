package pass

import (
	"fmt"
	"math/big"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/hvil/parser"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type TypeCheck struct {
	Result *errors.Collector
	parser.VisitorContext
	CurrentResultType parser.Type
}

// program
func (pass *TypeCheck) VisitProgram(program *parser.Program) {}

// packages
func (pass *TypeCheck) VisitPackage(pkg *parser.Package) {}

// functions
func (pass *TypeCheck) VisitFunction(function *parser.Function) {}

func (pass *TypeCheck) VisitFunctionVariableDeclaration(declaration *parser.FunctionVariableDeclaration) {
}

// basic blocks
func (pass *TypeCheck) VisitBlock(block *parser.BasicBlock) {}

func (pass *TypeCheck) VisitJump(terminator *parser.Jump) {}

func (pass *TypeCheck) VisitConditionalJump(terminator *parser.ConditionalJump) {
	_, ok := terminator.Condition.Type().(parser.ScalarType)
	if !ok {
		pass.Result.Err(
			terminator.Pos,
			"TypeError",
			fmt.Sprintf("only scalar values can be used as conditions in conditional jumps, but '%s' was given", terminator.Condition.Type().String()),
		)
	}
}

func (pass *TypeCheck) VisitReturn(terminator *parser.Return) {}

// instructions
func (pass *TypeCheck) VisitInstruction(instr *parser.Instruction) {
	pass.CurrentResultType = nil
}

func (pass *TypeCheck) VisitWriteRegister(write *parser.WriteRegister) {
	pass.CurrentResultType = write.Type
}

func (pass *TypeCheck) VisitWriteVariable(write *parser.WriteVariable) {
	pass.CurrentResultType = write.Declaration.Type()
}

func (pass *TypeCheck) VisitReadRegister(read *parser.ReadRegister) {}

func (pass *TypeCheck) VisitReadVariable(read *parser.ReadVariable) {}

// operations
func (pass *TypeCheck) VisitAluOperation(op *parser.AluOperation) {
	definition := aluDefinitions[op.Name]
	pass.checkSignatureType(op.Pos, definition, op.Args.Items)
}

func (pass *TypeCheck) VisitLocalCall(op *parser.LocalCall) {
	definition := op.Declaration.Type()
	pass.checkSignatureType(op.Pos, definition, op.Args.Items)
}

func (pass *TypeCheck) VisitDebugOperation(op *parser.DebugOperation) {
	definition := debugDefinitions[op.Name]
	pass.checkSignatureType(op.Pos, definition, op.Args.Items)
}

func (pass *TypeCheck) VisitScalarLiteral(op *parser.ScalarLiteral) {
	if pass.CurrentResultType == nil {
		pass.Result.Err(
			op.Pos,
			"TypeError",
			fmt.Sprintf("expected no return value but got a literal '%d'", op.Value),
		)
	} else {
		scalar, ok := pass.CurrentResultType.(parser.ScalarType)
		if !ok {
			pass.Result.Err(
				op.Pos,
				"TypeError",
				fmt.Sprintf("a literal cannot be written to '%s'", pass.CurrentResultType),
			)
		} else {
			two := big.NewInt(2)
			bitsize := big.NewInt(int64(scalar.BitSize))

			max := two.Exp(two, bitsize, nil)
			given := big.NewInt(int64(op.Value))

			if given.Cmp(max) >= 0 {
				pass.Result.Err(
					op.Pos,
					"TypeError",
					fmt.Sprintf("the literal value %d does not fit into %d bits", op.Value, scalar.BitSize),
				)
			}
		}
	}
}

func (pass *TypeCheck) checkSignatureType(
	position lexer.Position,
	definition parser.FunctionType,
	givenArgs []parser.Read,
) {
	if len(definition.Parameters.Items) != len(givenArgs) {
		pass.Result.Err(
			position,
			"TypeError",
			fmt.Sprintf("expected %d arg(s) but got %d", len(definition.Parameters.Items), len(givenArgs)),
		)
	} else {
		for i := range definition.Parameters.Items {
			expected := definition.Parameters.Items[i]
			actual := givenArgs[i]
			if !actual.Type().Equals(expected) {
				pass.Result.Err(
					actual.Position(),
					"TypeError",
					fmt.Sprintf("expected value of type '%s' but got '%s'", expected, actual.Type()),
				)
			}
		}
	}

	if pass.CurrentResultType != nil && definition.ReturnValue == nil {
		pass.Result.Err(
			position,
			"TypeError",
			fmt.Sprintf("expected return value of type '%s' none is provided", pass.CurrentResultType),
		)
	} else if pass.CurrentResultType == nil && definition.ReturnValue != nil {
		pass.Result.Err(
			position,
			"TypeError",
			fmt.Sprintf("this operation returns value of type '%s' but no return value is expected", definition.ReturnValue),
		)
	} else if pass.CurrentResultType != nil && definition.ReturnValue != nil && !pass.CurrentResultType.Equals(definition.ReturnValue) {
		pass.Result.Err(
			position,
			"TypeError",
			fmt.Sprintf("this operation returns value of type '%s' but '%s' is expected", definition.ReturnValue, pass.CurrentResultType),
		)
	}
}

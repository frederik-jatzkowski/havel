package parser_test

import (
	"bytes"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function/block/terminator"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types/scalar"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types/tuple"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/pass/parser"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		syntaxError bool
		expected    program.Program
	}{
		{
			name: "whitespace and comment",
			input: `
// hi
			`,
			syntaxError: true,
		},
		{
			name: "Ident",
			input: `
// hi asd
asd
			`,
			syntaxError: true,
		},
		{
			name: "single empty function",
			input: `
		func main() {
			block entry {} => return;
		}
					`,
			expected: program.Program{
				Functions: []*function.Function{
					{
						Name: "main",
						BasicBlocks: []*block.Block{
							{
								Ident:      "entry",
								Terminator: &terminator.Return{},
							},
						},
					},
				},
			},
		},
		{
			name: "nested tuples",
			input: `
func main (a1 : [8,16,[[32],16,8]]) {
	block entry {} => return;
}
			`,
			expected: program.Program{
				Functions: []*function.Function{
					{
						Name: "main",
						Parameters: tool.List[*stack.Decl]{
							Items: []*stack.Decl{
								{
									Name: "a1",
									DeclaredType: tuple.Type{
										Members: []types.Type{
											scalar.Type{BitSize: 8},
											scalar.Type{BitSize: 16},
											tuple.Type{
												Members: []types.Type{
													tuple.Type{
														Members: []types.Type{
															scalar.Type{BitSize: 32},
														},
													},
													scalar.Type{BitSize: 16},
													scalar.Type{BitSize: 8},
												},
											},
										},
									},
								},
							},
						},
						BasicBlocks: []*block.Block{
							{
								Ident:      "entry",
								Terminator: &terminator.Return{},
							},
						},
					},
				},
			},
		},
		{
			name: "function with return values",
			input: `
func main () => (r1 : [8, 16]) {
	block entry {} => return;
}
			`,
			expected: program.Program{
				Functions: []*function.Function{
					{
						Name: "main",
						ReturnValue: &stack.Decl{
							Name: "r1",
							DeclaredType: tuple.Type{
								Members: []types.Type{
									scalar.Type{BitSize: 8},
									scalar.Type{BitSize: 16},
								},
							},
						},
						BasicBlocks: []*block.Block{
							{
								Ident:      "entry",
								Terminator: &terminator.Return{},
							},
						},
					},
				},
			},
		},
		{
			name: "function with local variables",
			input: `
func main () {
	declare (
		a : [8,8],
		i1 : 16,
		i2 : 32
	);

	block entry {} => return;
}
			`,
			expected: program.Program{
				Functions: []*function.Function{
					{
						Name: "main",
						LocalDecls: tool.List[*stack.Decl]{
							Items: []*stack.Decl{
								{
									Name: "a",
									DeclaredType: tuple.Type{
										Members: []types.Type{
											scalar.Type{BitSize: 8},
											scalar.Type{BitSize: 8},
										},
									},
								},
								{
									Name:         "i1",
									DeclaredType: scalar.Type{BitSize: 16},
								},
								{
									Name:         "i2",
									DeclaredType: scalar.Type{BitSize: 32},
								},
							},
						},
						BasicBlocks: []*block.Block{
							{
								Ident:      "entry",
								Terminator: &terminator.Return{},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualProgram, err := parser.Parse(test.name, bytes.NewBufferString(test.input))
			if test.syntaxError != (err != nil) {
				t.Errorf("should error (%t), err: %s", test.syntaxError, err)
			}

			if test.syntaxError {
				return
			}

			if !reflect.DeepEqual(test.expected, actualProgram) {
				t.Errorf("expected '%+v' but got '%+v'", test.expected, actualProgram)
			}
		})
	}
}

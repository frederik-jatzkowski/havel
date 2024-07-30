package parser_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/frederik-jatzkowski/havel/internal/hvil/parser"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		syntaxError bool
		expected    parser.Package
	}{
		{
			name: "whitespace and comment",
			input: `
// hi
			`,
			syntaxError: true,
		},
		{
			name: "identifier",
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
			expected: parser.Package{
				Functions: []*parser.Function{
					{
						Name: "main",
						BasicBlocks: []*parser.BasicBlock{
							{
								Identifier: "entry",
								Terminator: &parser.Return{},
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
			expected: parser.Package{
				Functions: []*parser.Function{
					{
						Name: "main",
						Parameters: parser.CommaSeparatedList[*parser.FunctionVariableDeclaration]{
							Items: []*parser.FunctionVariableDeclaration{
								{
									Name: "a1",
									DeclaredType: parser.TupleType{
										Members: []parser.Type{
											parser.ScalarType{BitSize: 8},
											parser.ScalarType{BitSize: 16},
											parser.TupleType{
												Members: []parser.Type{
													parser.TupleType{
														Members: []parser.Type{
															parser.ScalarType{BitSize: 32},
														},
													},
													parser.ScalarType{BitSize: 16},
													parser.ScalarType{BitSize: 8},
												},
											},
										},
									},
								},
							},
						},
						BasicBlocks: []*parser.BasicBlock{
							{
								Identifier: "entry",
								Terminator: &parser.Return{},
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
			expected: parser.Package{
				Functions: []*parser.Function{
					{
						Name: "main",
						ReturnValue: &parser.FunctionVariableDeclaration{
							Name: "r1",
							DeclaredType: parser.TupleType{
								Members: []parser.Type{
									parser.ScalarType{BitSize: 8},
									parser.ScalarType{BitSize: 16},
								},
							},
						},
						BasicBlocks: []*parser.BasicBlock{
							{
								Identifier: "entry",
								Terminator: &parser.Return{},
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
			expected: parser.Package{
				Functions: []*parser.Function{
					{
						Name: "main",
						LocalDeclarations: parser.CommaSeparatedList[*parser.FunctionVariableDeclaration]{
							Items: []*parser.FunctionVariableDeclaration{
								{
									Name: "a",
									DeclaredType: parser.TupleType{
										Members: []parser.Type{
											parser.ScalarType{BitSize: 8},
											parser.ScalarType{BitSize: 8},
										},
									},
								},
								{
									Name:         "i1",
									DeclaredType: parser.ScalarType{BitSize: 16},
								},
								{
									Name:         "i2",
									DeclaredType: parser.ScalarType{BitSize: 32},
								},
							},
						},
						BasicBlocks: []*parser.BasicBlock{
							{
								Identifier: "entry",
								Terminator: &parser.Return{},
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

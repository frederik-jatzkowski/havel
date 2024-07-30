package pass

import "github.com/frederik-jatzkowski/havel/internal/hvil/parser"

var (
	debugDefinitions = map[string]parser.FunctionType{
		"print_u_32": {
			Parameters: parser.CommaSeparatedList[parser.Type]{
				Items: []parser.Type{
					parser.ScalarType{BitSize: 32},
				},
			},
		},
	}
	aluDefinitions = map[string]parser.FunctionType{
		"add_u_32": {
			Parameters: parser.CommaSeparatedList[parser.Type]{
				Items: []parser.Type{
					parser.ScalarType{BitSize: 32},
					parser.ScalarType{BitSize: 32},
				},
			},
			ReturnValue: parser.ScalarType{BitSize: 32},
		},
		"sub_u_32": {
			Parameters: parser.CommaSeparatedList[parser.Type]{
				Items: []parser.Type{
					parser.ScalarType{BitSize: 32},
					parser.ScalarType{BitSize: 32},
				},
			},
			ReturnValue: parser.ScalarType{BitSize: 32},
		},
		"lt_u_32": {
			Parameters: parser.CommaSeparatedList[parser.Type]{
				Items: []parser.Type{
					parser.ScalarType{BitSize: 32},
					parser.ScalarType{BitSize: 32},
				},
			},
			ReturnValue: parser.ScalarType{BitSize: 32},
		},
	}
)

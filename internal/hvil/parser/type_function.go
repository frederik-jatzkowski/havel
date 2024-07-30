package parser

type FunctionType struct {
	Parameters  CommaSeparatedList[Type] `parser:"'func' '(' @@ ')'"`
	ReturnValue Type                     `parser:"( '=>' @@ )?"`
}

func (t FunctionType) String() string {
	result := "func("
	for _, param := range t.Parameters.Items {
		result += param.String()
	}
	result += ")"

	if t.ReturnValue != nil {
		result += "=>(" + t.ReturnValue.String() + ")"
	}

	return result
}

func (t FunctionType) Equals(other Type) bool {
	return t.String() == other.String()
}

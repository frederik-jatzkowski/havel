package types

import (
	"fmt"
	"strings"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Composite struct {
	tool.Node[Composite]

	Fields tool.List[Type] `parser:"'{' @@ '}'"`
}

func (node *Composite) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *Composite) EqualsDetailed(other Type) error {
	composite, ok := other.(*Composite)
	if !ok {
		return node.Errorf("other type is not a composite")
	}

	if len(node.Fields.Items) != len(composite.Fields.Items) {
		return fmt.Errorf("composite count mismatch: expected %d, got %d", len(node.Fields.Items), len(composite.Fields.Items))
	}

	for i, field := range node.Fields.Items {
		otherField := composite.Fields.Items[i]
		if err := field.EqualsDetailed(otherField); err != nil {
			return err
		}
	}

	return nil
}

func (node *Composite) Bytes() int {
	total := 0
	for _, field := range node.Fields.Items {
		total += field.Bytes()
	}

	return total
}

func (node *Composite) CanDoArithmetics() bool {
	return false
}

func (node *Composite) Dereference(fields []uint) (t Type, offset uint, err error) {
	if len(fields) == 0 {
		return node, 0, nil
	}

	if fields[0] > uint(len(node.Fields.Items)-1) {
		return node, 0, fmt.Errorf("not enough fields: has %d but wants %d", len(node.Fields.Items), fields[0])
	}

	for i, field := range node.Fields.Items {
		if uint(i) == fields[0] {
			result, innerOffset, err := field.Dereference(fields[1:])

			return result, innerOffset + offset, err
		}

		offset += uint(field.Bytes())
	}

	return node, offset, nil
}

func (node *Composite) String() string {
	result := strings.Builder{}

	result.WriteString("{ ")

	for i, item := range node.Fields.Items {
		result.WriteString(item.String())

		if i < len(node.Fields.Items)-1 {
			result.WriteString(", ")
		}
	}

	result.WriteString(" }")

	return result.String()
}

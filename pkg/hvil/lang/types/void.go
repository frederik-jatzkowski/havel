package types

import "github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"

type Void struct {
	tool.Node[Void]
}

func (node *Void) String() string {
	return "void"
}

func (node *Void) CanBeAssigned(other Type) bool {
	return node.Equals(other)
}

func (node *Void) Equals(other Type) bool {
	return node.EqualsDetailed(other) == nil
}

func (node *Void) EqualsDetailed(other Type) error {
	_, ok := other.(*Void)
	if !ok {
		return node.Errorf("%s is not void", other)
	}

	return nil
}

func (node *Void) Bytes() int {
	return 0
}

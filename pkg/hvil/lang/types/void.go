package types

type Void struct{}

func (node Void) String() string {
	return "void"
}

func (node Void) CanBeAssigned(_ Type) bool {
	return true
}

func (node Void) Bytes() int {
	return 0
}

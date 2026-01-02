package types

type Void struct{}

func (node Void) String() string {
	return "void"
}

func (node Void) CanBeAssigned(_ Type) bool {
	return true
}

func (node Void) Equals(other Type) bool {
	_, ok := other.(Void)

	return ok
}

func (node Void) Bytes() int {
	return 0
}

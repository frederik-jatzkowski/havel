package types

type Void struct{}

func (v Void) String() string {
	return "void"
}

func (v Void) CanBeAssigned(_ Type) bool {
	return true
}

func (v Void) Bytes() int {
	return 0
}

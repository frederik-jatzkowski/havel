package names

//go:generate go tool go-enum ./kind.go

// Kind represents the kind of the scope
// ENUM(function, block, variable, register)
type Kind byte

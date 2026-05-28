package names

//go:generate go tool go-enum ./kind.go

// Kind represents the kind of the scope
// ENUM(identifier, struct_member)
type Kind byte

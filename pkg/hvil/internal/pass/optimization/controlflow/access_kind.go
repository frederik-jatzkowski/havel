package controlflow

//go:generate go tool go-enum ./access_kind.go

// AccessKind indicates if an access is a read or a write
// ENUM(READ, WRITE)
type AccessKind uint8

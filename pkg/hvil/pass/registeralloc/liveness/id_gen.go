package liveness

type InstructionID uint

type Liveness[T any] struct {
	LivenessPass T `parser:"" json:",omitempty"`
}

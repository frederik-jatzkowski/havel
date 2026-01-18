package architecture

type Architecture interface {
	GeneralPurposeRegisters() []Register
	ScratchRegisters() []Register
}

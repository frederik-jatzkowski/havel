package architecture

type Architecture interface {
	GetArgRegister() (Register, bool)
	GetResultRegister() (Register, bool)
	GetGeneralPurposeRegister() (Register, bool)
	ReturnGeneralPurposeRegisters(r ...Register)
	GetScratchRegister() (Register, bool)
	ReturnScratchRegisters(r ...Register)
}

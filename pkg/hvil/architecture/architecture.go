package architecture

type Architecture interface {
	GetArgRegister() (Register, bool)
	GetResultRegister() (Register, bool)
	GetGeneralPurposeRegister() (Register, bool)
	GetScratchRegister() (Register, bool)
}

package fibonacci

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/frederik-jatzkowski/havel/pkg/hvil"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine"
)

func BenchmarkCollatz(b *testing.B) {
	n := 214

	fmt.Printf("n = 0b%b\n", n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SolveCollatz(214)
	}
}

func SolveCollatz(n uint64) uint64 {
	if n == 1 {
		return 1
	}

	if n%2 == 0 {
		return SolveCollatz(n/2) + 1
	}

	return SolveCollatz(3*n+1) + 1
}

func BenchmarkCollatz_virtualmachine(b *testing.B) {
	filePath := "./src.hvil"
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	compiler := hvil.NewCompiler()

	program, err := compiler.Compile(filePath, file)
	if err != nil {
		panic(fmt.Errorf("compilation failed:\n %w", err))
	}

	asm, err := program.GenerateVirtualMachineAssembly()
	cobra.CheckErr(err)

	bc, err := asm.Assemble()
	cobra.CheckErr(err)

	vm := virtualmachine.New(
		1024*1024,
		os.Stdin,
		io.Discard,
		io.Discard,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = vm.Execute(bc)
		if err != nil {
			panic(errors.Join(
				errors.New("runtime error"),
				err,
			))
		}
	}
}

func BenchmarkCollatz_interpreted(b *testing.B) {
	filePath := "./src.hvil"
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	compiler := hvil.NewCompiler()

	program, err := compiler.Compile(filePath, file)
	if err != nil {
		panic(fmt.Errorf("compilation failed:\n %w", err))
	}

	vm := runtime.New(
		1024*1024,
		os.Stdin,
		io.Discard,
		io.Discard,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = program.Execute(vm)
		if err != nil {
			panic(errors.Join(
				errors.New("runtime error"),
				err,
			))
		}
	}
}

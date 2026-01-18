package hvil

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine"
)

func TestCompiler(t *testing.T) {
	type ExpectedError struct {
		Contains string `json:"contains"`
	}
	type ExpectedOutput struct {
		Compiler struct {
			Errors []ExpectedError `json:"errors"`
		} `json:"compiler"`
		Execution struct {
			StdoutLines []string        `json:"stdout_lines"`
			Errors      []ExpectedError `json:"errors"`
		} `json:"execution"`
	}

	err := filepath.WalkDir("./testdata", func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsDir() {
			return nil
		}

		srcPath := filepath.Join(path, "src.hvil")
		src, err := os.ReadFile(srcPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}

			return err
		}

		t.Run(path, func(t *testing.T) {
			t.Parallel()

			expectedErrsData, err := os.ReadFile(filepath.Join(path, "spec.json"))
			require.NoError(t, err)

			var expectedOutput ExpectedOutput

			err = json.Unmarshal(expectedErrsData, &expectedOutput)
			require.NoError(t, err)

			compiler := NewCompiler()
			program, err := compiler.Compile(srcPath, bytes.NewBuffer(src))

			for _, expectedErr := range expectedOutput.Compiler.Errors {
				assert.ErrorContains(t, err, expectedErr.Contains)
			}

			if len(expectedOutput.Compiler.Errors) == 0 && err != nil {
				t.Error("unexpected compiler error")
			}

			if err != nil {
				t.Logf("found errors:\n%s", err)
				return
			}

			stdout := bytes.NewBuffer(nil)

			vm := runtime.New(1024, bytes.NewBuffer(nil), stdout, io.Discard)

			err = program.Execute(vm)

			if len(expectedOutput.Execution.Errors) == 0 && err != nil {
				require.NoError(t, err)
			}

			if len(expectedOutput.Execution.Errors) > 0 && err == nil {
				require.Error(t, err)
			}

			for _, expectedErr := range expectedOutput.Execution.Errors {
				assert.ErrorContains(t, err, expectedErr.Contains)
			}

			actualLines := strings.Split(stdout.String(), "\n")
			for i, expectedLine := range expectedOutput.Execution.StdoutLines {
				if i > len(actualLines)-1 {
					t.Errorf("expected %d lines but got %d", len(expectedOutput.Execution.StdoutLines), len(actualLines))
					continue
				}

				assert.Contains(t, actualLines[i], expectedLine)
			}
		})

		return nil
	})

	require.NoError(t, err)
}

func TestCompiler_OnVirtualMachine(t *testing.T) {
	type ExpectedError struct {
		Contains string `json:"contains"`
	}
	type ExpectedOutput struct {
		Compiler struct {
			Errors []ExpectedError `json:"errors"`
		} `json:"compiler"`
		Execution struct {
			StdoutLines []string        `json:"stdout_lines"`
			Errors      []ExpectedError `json:"errors"`
		} `json:"execution"`
	}

	err := filepath.WalkDir("./testdata", func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsDir() {
			return nil
		}

		srcPath := filepath.Join(path, "src.hvil")
		src, err := os.ReadFile(srcPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}

			return err
		}

		t.Run(path, func(t *testing.T) {
			t.Parallel()

			require.NotPanics(t, func() {
				expectedErrsData, err := os.ReadFile(filepath.Join(path, "spec.json"))
				require.NoError(t, err)

				var expectedOutput ExpectedOutput

				err = json.Unmarshal(expectedErrsData, &expectedOutput)
				require.NoError(t, err)

				t.Log("Source:\n" + string(src))

				compiler := NewCompiler()
				program, err := compiler.Compile(srcPath, bytes.NewBuffer(src))

				for _, expectedErr := range expectedOutput.Compiler.Errors {
					assert.ErrorContains(t, err, expectedErr.Contains)
				}

				if len(expectedOutput.Compiler.Errors) == 0 && err != nil {
					t.Error("unexpected compiler error")
				}

				if err != nil {
					t.Logf("found errors:\n%s", err)
					return
				}

				asm, err := program.GenerateVirtualMachineAssembly()
				require.NoError(t, err)

				t.Log("Assembly:\n" + asm.String())

				byteCode, err := asm.Assemble()
				require.NoError(t, err)

				t.Log("Bytecode:\n" + byteCode.String())

				stdout := bytes.NewBuffer(nil)
				vm := virtualmachine.New(1024*1024, bytes.NewBuffer(nil), stdout, io.Discard)

				err = vm.Execute(byteCode)

				if len(expectedOutput.Execution.Errors) == 0 && err != nil {
					require.NoError(t, err)
				}

				if len(expectedOutput.Execution.Errors) > 0 && err == nil {
					require.Error(t, err)
				}

				for _, expectedErr := range expectedOutput.Execution.Errors {
					assert.ErrorContains(t, err, expectedErr.Contains)
				}

				actualLines := strings.Split(stdout.String(), "\n")
				for i, expectedLine := range expectedOutput.Execution.StdoutLines {
					if i > len(actualLines)-1 {
						t.Errorf("expected %d lines but got %d", len(expectedOutput.Execution.StdoutLines), len(actualLines))
						continue
					}

					assert.Contains(t, actualLines[i], expectedLine)
				}
			})
		})

		return nil
	})

	require.NoError(t, err)
}

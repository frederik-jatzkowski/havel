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
)

func TestCompiler(t *testing.T) {
	type ExpectedCompilerError struct {
		Contains string `json:"contains"`
	}
	type ExpectedOutput struct {
		Compiler struct {
			Errors []ExpectedCompilerError `json:"errors"`
		} `json:"compiler"`
		Execution struct {
			StdoutLines []string `json:"stdout_lines"`
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
			require.NoError(t, err)

			actualLines := strings.Split(stdout.String(), "\n")
			for i, actualLine := range actualLines {
				if i == len(actualLines)-1 {
					continue
				}

				assert.Contains(t, actualLine, expectedOutput.Execution.StdoutLines[i])
			}
		})

		return nil
	})

	require.NoError(t, err)
}

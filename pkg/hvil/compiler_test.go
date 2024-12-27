package hvil

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestCompiler(t *testing.T) {
	type ExpectedCompilerError struct {
		Message string `json:"message"`
	}
	type ExpectedOutput struct {
		Compiler struct {
			Errors []ExpectedCompilerError `json:"errors"`
		} `json:"compiler"`
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

			expectedErrsData, err := os.ReadFile(filepath.Join(path, "output.json"))
			require.NoError(t, err)

			var expectedOutput ExpectedOutput

			err = json.Unmarshal(expectedErrsData, &expectedOutput)
			require.NoError(t, err)

			expectedErrsMap := make(map[string]ExpectedCompilerError, len(expectedOutput.Compiler.Errors))
			for _, expectedErr := range expectedOutput.Compiler.Errors {
				expectedErrsMap[expectedErr.Message] = expectedErr
			}

			compiler := NewCompiler()
			_, actualErrs := compiler.Compile(srcPath, bytes.NewBuffer(src))

			for _, actualErr := range actualErrs {
				errMsg := actualErr.Error()
				_, exists := expectedErrsMap[errMsg]
				assert.Truef(t, exists, "error message '%s' was unexpected", errMsg)

				if exists {
					delete(expectedErrsMap, errMsg)
				}
			}

			for _, expectedErr := range expectedErrsMap {
				t.Fail()
				t.Logf("remaining error expectation: '%s'", expectedErr.Message)
			}
		})

		return nil
	})

	require.NoError(t, err)
}

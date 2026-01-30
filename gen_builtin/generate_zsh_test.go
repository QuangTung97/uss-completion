package gen_builtin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateZshBuiltinCompletionScript(t *testing.T) {
	const outputPath = "output/output_zsh.sh"
	GenerateZshBuiltinCompletionScript(
		outputPath,
		Command{Name: "cd", NoSpace: true},
		Command{Name: "cat"},
		Command{Name: "ls"},
	)

	outputData, err := os.ReadFile(outputPath)
	require.Equal(t, nil, err)

	expectedData, err := os.ReadFile("../zsh_builtin_complete.sh")
	require.Equal(t, nil, err)

	assert.Equal(t, string(expectedData), string(outputData))
}

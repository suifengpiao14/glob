package glob

import (
	"embed"
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed example
var f embed.FS

func TestGlobFS(t *testing.T) {
	fileList, err := GlobFS(f, "example/**/*.xml")
	require.NoError(t, err)
	fmt.Println(fileList)
}

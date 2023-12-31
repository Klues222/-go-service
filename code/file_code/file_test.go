package file_code

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	f, err := os.Open("testdata/my_file.txt")
	require.NoError(t, err)
	data := make([]byte, 64)
	n, err := f.Read(data)
	fmt.Println(n)
	require.Error(t, err)

	n, err = f.WriteString("Hello")
	fmt.Println(n)
	//bad file descriptor 不可写
	fmt.Println(err)
	require.Error(t, err)
	f.Close()

	f, err = os.OpenFile("testdata/my_file.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	require.NoError(t, err)
	n, err = f.WriteString("world")
	fmt.Println(n)
	require.NoError(t, err)

	f, err = os.Create("testdata/my_file_copy.txt")
	require.NoError(t, err)
	n, err = f.WriteString("hello world")
	fmt.Println(n)
	require.NoError(t, err)
}

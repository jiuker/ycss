package filePath

import (
	"fmt"
	"testing"
)

func TestNewFilePath(t *testing.T) {
	f := NewFilePath("./res/sample/rn.js")
	fmt.Println(f.GetFilePath())
	fmt.Println(f.GetFileDir())
	fmt.Println(f.GetFileName())
	fmt.Println(f.GetFileType())
}

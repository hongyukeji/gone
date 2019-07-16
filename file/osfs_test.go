package file_test

import (
	"iki-go/file"
	"testing"
)

func TestOsFs_Create(t *testing.T) {
	fs := file.NewOsFs()
	f, err := fs.Create("hello world.log")
	defer f.Close()
	if err != nil {
		t.Log(f.Name())
	}
	f.WriteString("hello world!!")

}

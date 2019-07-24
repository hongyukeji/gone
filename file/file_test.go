package file_test

import (
	"os"
	"testing"

	"github.com/wx11055/gone/file"
)

const fileName = "./file_test.go"

func TestReadFileWithIoutil(t *testing.T) {
	err := file.ReadFileWithIoutil(fileName, func(line int, content string) {
		t.Logf("line: %d , content: %s", line, content)
	})
	t.Log(err)
}
func TestReadLineWithOs(t *testing.T) {
	err := file.ReadLineWithOs(fileName, func(line int, content string) {
		t.Logf("line: %d , content: %s", line, content)
	})
	t.Log(err)
}
func TestReadLineWithFile(t *testing.T) {
	err := file.ReadLineWithFile(fileName, func(line int, content string) {
		t.Logf("line: %d , content: %s", line, content)
	})
	t.Log(err)
}

func TestReadLineWithBuffer(t *testing.T) {
	err := file.ReadLineWithBufio(fileName, func(line int, content string) {
		t.Logf("line: %d , content: %s", line, content)
	})
	t.Log(err)
}

const writeFileName = "file2_test.go"
const writeFileBody = `
package file_test

import (
	"github.com/wx11055/gone/file"
	"testing"
)

func TestWriteFileWithIoutil2(t *testing.T) {
	if err := file.WriteFileWithIoutil(writeFileName, []byte(writeFileBody)); err != nil {
		t.Errorf("WriteFileWithIoutil() error = %v", err)
	}
}
`

func TestWriteFileWithIoutil(t *testing.T) {
	if err := file.WriteFileWithIoutil(writeFileName, []byte(writeFileBody)); err != nil {
		t.Errorf("WriteFileWithIoutil() error = %v", err)
	}
}
func TestWriteFileWithOS(t *testing.T) {
	if err := file.WriteFileWithOS(writeFileName, []byte(writeFileBody)); err != nil {
		t.Errorf("WriteFileWithOS() error = %v", err)
	}
}

func TestWriteFileBufio(t *testing.T) {
	if err := file.WriteFileBufio(writeFileName, os.O_RDWR|os.O_CREATE, []byte(writeFileBody)); err != nil {
		t.Errorf("WriteFileBufio() error = %v", err)
	}
}

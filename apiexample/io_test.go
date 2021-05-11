package apiexample

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func ReadFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}

func TestIoPack(t *testing.T) {
	if from, err := ReadFrom(os.Stdin, 0); err != nil {
		fmt.Println(from)
	}
}

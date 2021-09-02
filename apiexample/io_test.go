package apiexample

import (
	"fmt"
	"io"
	"os"
	"strconv"
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

func TestDemo(t *testing.T) {
	rawHex := "39A7F8"
	println("001110011010011111111000")
	i, err := strconv.ParseUint(rawHex, 16, 32)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%024b\n", i)

	binS := "1001101110011110110101"
	i1, err1 := strconv.ParseUint(binS, 2, 64)
	println("26E7B5")
	if err1 != nil {
		fmt.Printf("%s", err1)
	}
	fmt.Printf("%x\n", i1)
}

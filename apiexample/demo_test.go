package apiexample

import (
	"fmt"
	"testing"
)

func TestDemo1(t *testing.T) {
	fmt.Println("test 123 ")
	a := "123"
	t.Logf("这是一条打因日志:%s", a)
}

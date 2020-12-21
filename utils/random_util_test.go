package utils

import "testing"

func TestGetRandomNum(t *testing.T) {
	num := GetRandomNum(6)
	t.Log(num)
}

package services

import (
	"strings"
	"testing"
)

func TestSendCode(t *testing.T) {
	code, err := SendCode("13260407984", 0)
	if err != nil {
		t.Log(code)
	}
}

func TestSmsBuilder(t *testing.T) {
	builder, _ := smsBuilder("123456")
	t.Log(builder)
}

func TestCheckCode(t *testing.T) {
	code, err := CheckCode("13210242064", "123456")
	t.Log(code, err)
}

func TestCheckCode2(t *testing.T) {
	var str = "    6000   "
	trim := strings.TrimSpace(str)
	t.Log(len(trim))
}

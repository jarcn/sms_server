package models

import (
	"testing"
	"time"
)

func TestSmsCodeBean_Insert(t *testing.T) {
	bean := TSmsRecord{
		PhoneNum: "13210242096",
		SendTime: time.Now().Unix(),
		Type:     0,
		SmsCode:  "1234",
	}
	if insert, err := bean.Insert(); err != nil {
		t.Log(err)
	} else {
		t.Log(insert)
	}
}

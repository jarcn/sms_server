package utils

import (
	"testing"
	"time"
)

func TestExpiredMap_Get(t *testing.T) {
	expiredMap := NewExpiredMap()
	expiredMap.Set("key", "value", 5)
	time.Sleep(time.Second * 6)
	if b, value := expiredMap.Get("key"); b {
		t.Log(value)
	} else {
		t.Log("time out")
	}
}

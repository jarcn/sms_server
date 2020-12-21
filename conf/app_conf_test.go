package conf

import (
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	timeS := 15 * time.Minute
	t.Log(timeS.Seconds())
	var timeOut string
	if "" == timeOut {
		timeOut = "100"
	} else {
		timeOut = "900"
	}
	t.Log(timeOut)
}

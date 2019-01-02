package profile

import (
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	defer Duration(time.Now(),"Test Duration")
	time.Sleep(10*time.Second)
}


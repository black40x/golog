package golog

import (
	"testing"
	"time"
)

func Test_IsCurrentDate(t *testing.T) {
	if !isCurrentDay(time.Now()) {
		t.Error("Current day time.Now() wanted true")
	}
}

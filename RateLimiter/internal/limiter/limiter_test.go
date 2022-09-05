package limiter

import (
	"testing"
	"time"
)

func TestCheckSkip(t *testing.T) {
	limiter := NewLimiter(2, time.Hour)

	request_1 := limiter.CheckSkip()
	if request_1 {
		t.Fatal("error", request_1)
	}

	request_2 := limiter.CheckSkip()
	if request_2 {
		t.Fatal("error", request_2)
	}

	request_3 := limiter.CheckSkip()
	if !request_3 {
		t.Fatal("error", request_3)
	}
}

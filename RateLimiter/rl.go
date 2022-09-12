package main

import (
	"fmt"
	"log"
	"net/http"
	"ratelimiter/internal/limiter"
	"time"
)

type MyHandler struct {
	limiter *limiter.Limiter
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	skip := h.limiter.CheckSkip()

	if skip == true {
		fmt.Fprintf(w, "count is %s\n", "vo 1 ya v drugom gorode")
	} else {
		//
		fmt.Fprintf(w, "count is %s\n", "answer")
	}
}

func main() {
	l := limiter.NewLimiter(1, time.Second*10)
	h := &MyHandler{l}

	http.Handle("/count", h)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

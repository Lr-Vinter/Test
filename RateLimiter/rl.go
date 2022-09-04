package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type MyHandler struct {
	limiter *Limiter
}

type Limiter struct {
	counter chan int
	answer  chan bool

	maxRequest int
	timeline   time.Duration

	endtime chan bool
}

func (l *Limiter) sendTime() {
	l.counter <- 0
	interval := time.Tick(l.timeline)
	for range interval {
		l.endtime <- true
	}
}

func (l *Limiter) restartClock() {
	for range l.endtime {
		<-l.endtime
		l.counter <- 0
	}
}

func (l *Limiter) CheckSkip() {
	val := <-l.counter
	val++
	if val > l.maxRequest {
		l.answer <- true
	} else {
		l.answer <- false
	}
	l.counter <- val
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go h.limiter.CheckSkip()

	for val := range h.limiter.answer {
		if val == true {
			fmt.Fprintf(w, "count is %s\n", "vo 1 ya v drugom gorode")
		} else {
			fmt.Fprintf(w, "count is %s\n", "answer")
			//h.limiter.answer <- false //(? WHY)
		}
		return // ?
	}
}

func main() {

	l := &Limiter{make(chan int), make(chan bool), 2, time.Second * 5, make(chan bool)} // 2 requests by 5 sec

	go l.sendTime()
	go l.restartClock()

	h := &MyHandler{l}

	http.Handle("/count", h)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

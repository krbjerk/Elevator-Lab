package main

import (
	"fmt"
	"time"
)

// Timer struct to hold timer state
type Timer struct {
	endTime time.Time
	active  bool
}

// Declare timer object
var timer Timer

// Start() - Starts the timer for a given duration in seconds
func (t *Timer) Start(duration float64) {
	t.endTime = time.Now().Add(time.Duration(duration * float64(time.Second)))
	t.active = true
	fmt.Println("Timer started.")
}

// Stop() - Stops the timer
func (t *Timer) Stop() {
	t.active = false
	fmt.Println("Timer stopped.")
}

// TimedOut() - Checks if the timer has expired
func (t *Timer) TimedOut() bool {

	if t.active && time.Now().After(t.endTime) {
		t.Stop()
		return true
	} else {
		return false
	}
}

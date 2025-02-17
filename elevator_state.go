package main

import (
	"Driver-go/elevio"
)

type elevator struct {
	floor    int
	dirn     elevio.MotorDirection
	request  [4][3]int
	behavior ElevatorBehavior
}

type ElevatorBehavior int

const (
	EB_idle     ElevatorBehavior = 0
	EB_DoorOpen                  = 1
	EB_Moving                    = 2
)

// Create elevator object

// Create state machine functions:
// Needs to react on button press and floor

func elevatorInit() {

}

func buttonPress() {

}

func floorArrival() {

}

func doorTimeout() {

}

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

// Created a button type for use in statemachine to decide further action.
type Button int

const (
	D_Down Button = -1
	D_Stop        = 0
	D_Up          = 1
)

// Create elevator object
var Elevator elevator

// Create state machine functions:
// Needs to react on button press and floor

func elevatorInit() {
	// TODO: Make sure this is really necessary.
	elevio.SetMotorDirection(elevio.MD_Down)
	Elevator.dirn = elevio.MD_Down
	Elevator.behavior = EB_Moving

}

func buttonPress(btnFloor int, btnType Button) {
	switch Elevator.behavior {
	case EB_DoorOpen:

	case EB_Moving:

	case EB_idle:
	}
}

func floorArrival() {
	switch Elevator.behavior {
	case EB_Moving:
	}
}

func doorTimeout() {
	switch Elevator.behavior {
	case EB_DoorOpen:
	}
}

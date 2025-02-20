package main

import (
	"Driver-go/elevio"
)

type elevator struct {
	floor    int
	dirn     elevio.MotorDirection
	requests [4][3]int
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
	// Requests can happen at any point.

	switch Elevator.behavior {
	case EB_DoorOpen:
		// Clear if request is for same floor
		// Otherwise request will be approved and stored for later
		if Elevator.floor == btnFloor {
			Elevator.floor = btnFloor
			// Start timer for door open
		} else {
			Elevator.requests[btnFloor][btnType] = 1
		}

	case EB_Moving:
		// Take request and store it for later
		Elevator.requests[btnFloor][btnType] = 1
	case EB_idle:
		// Take request and act on it immediately
		Elevator.requests[btnFloor][btnType] = 1
		// TODO

		// Based on the requests create a "twin" that knows what the elevator should do
		// We will use that twin to actually make the elevator do

	}
}

func floorArrival() {
	// Moving given implicitly

	// Check if the current request needs to stop at this floor
	//	Stop the cab
	//	Do something with the lights
	//	Clear request
	//	Open door and start timer

}

func doorTimeout() {
	// Idle given implicitly

	// check if there are requests
	// Act on them

}

// TODO: We need some request prioritization

// We need some system to handle requests

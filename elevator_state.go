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
	EB_Idle     ElevatorBehavior = 0
	EB_DoorOpen                  = 1
	EB_Moving                    = 2
)

// Created a button type for use in statemachine to decide further action.
type Direction int

const (
	D_Down Direction = -1
	D_Stop           = 0
	D_Up             = 1
)

type Button int

const (
	B_HallUp   Button = 0
	B_HallDown        = 1
	B_Cab             = 2
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
			// Start timer for door open
		} else {
			Elevator.requests[btnFloor][btnType] = 1
		}

	case EB_Moving:
		// Take request and store it for later
		Elevator.requests[btnFloor][btnType] = 1
	case EB_Idle:
		// Take request and act on it immediately
		Elevator.requests[btnFloor][btnType] = 1
		// TODO
		twin := Elevator.requestDirection()
		Elevator.dirn = twin.dirn
		Elevator.behavior = twin.behavior

		switch twin.behavior {
		case EB_DoorOpen:
			// Do something about the doorlight
			// Start the door timer
			Elevator = Elevator.requestsClearAtCurrentFloor()

		case EB_Moving:
			//
			elevio.SetMotorDirection(Elevator.dirn)

		case EB_Idle:
		}

		// SetAlllights(Elevator)
		// Print out something about the elevator having a new state.

	}
}

func floorArrival(newFloor int) {
	// Moving given implicitly

	// Check if the current request needs to stop at this floor
	//	Stop the cab
	//	Do something with the lights
	//	Clear request
	//	Open door and start timer

	// Print out that the elevator has arrived at a new floor
	Elevator.floor = newFloor
	elevio.SetFloorIndicator(Elevator.floor)

	switch Elevator.behavior {
	case EB_Moving:
		if Elevator.requestsShouldStop() {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			Elevator = Elevator.requestsClearAtCurrentFloor()
			// Start open door timer
			// Set all lights
			Elevator.behavior = EB_DoorOpen
		}
	default:

	}
	// Print the new state of elevator.
}

func doorTimeout() {
	// Idle given implicitly

	// check if there are requests
	// Act on them

	// Elevator printout
	switch Elevator.behavior {
	case EB_DoorOpen:
		twin := Elevator.requestDirection()
		Elevator.dirn = twin.dirn
		Elevator.behavior = twin.behavior

		switch Elevator.behavior {
		case EB_DoorOpen:
			// Start timer open door
			Elevator = Elevator.requestsClearAtCurrentFloor()
			// Set all lights
		case EB_Moving:

		case EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(Elevator.dirn)
		}

	default:

	}
	// Print elevator state
}

// TODO: We need some request prioritization

// We need some system to handle requests

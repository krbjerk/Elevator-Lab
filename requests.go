package main

//

import (
	"Driver-go/elevio"
)

type Twin struct {
	dirn     elevio.MotorDirection
	behavior ElevatorBehavior
}

// Functions for simplification

func (e elevator) requestAbove() bool {
	for i := e.floor + 1; i < 4; i++ {
		for j := 0; j < 3; j++ {
			if e.requests[i][j] == 1 {
				return true
			}
		}
	}
	return false
}

func (e elevator) requestBelow() bool {
	for i := 0; i < e.floor; i++ {
		for j := 0; j < 3; j++ {
			if e.requests[i][j] == 1 {
				return true
			}
		}
	}
	return false

}

func (e elevator) requestHere() bool {
	for i := 0; i < 3; i++ {
		if e.requests[e.floor][i] == 1 {
			return true
		}
	}
	return false
}

// request choose direction

func (e elevator) requestDirection() Twin {
	switch e.dirn {
	case elevio.MD_Up:
		if e.requestAbove() {
			return Twin{elevio.MD_Up, EB_Moving}
		} else if e.requestHere() {
			return Twin{elevio.MD_Down, EB_DoorOpen}
		} else if e.requestBelow() {
			return Twin{elevio.MD_Down, EB_Moving}
		} else {
			return Twin{elevio.MD_Stop, EB_Idle}
		}

	case elevio.MD_Down:
		if e.requestBelow() {
			return Twin{elevio.MD_Down, EB_Moving}
		} else if e.requestHere() {
			return Twin{elevio.MD_Up, EB_DoorOpen}
		} else if e.requestAbove() {
			return Twin{elevio.MD_Up, EB_Moving}
		} else {
			return Twin{elevio.MD_Stop, EB_Idle}
		}

	case elevio.MD_Stop:
		if e.requestHere() {
			return Twin{elevio.MD_Stop, EB_DoorOpen}
		} else if e.requestAbove() {
			return Twin{elevio.MD_Up, EB_Moving}
		} else if e.requestBelow() {
			return Twin{elevio.MD_Down, EB_Moving}
		} else {
			return Twin{elevio.MD_Stop, EB_Idle}
		}
	default:
		return Twin{elevio.MD_Stop, EB_Idle} // Must include default to not get no return error.
	}
}

// request should stop
func (e elevator) requestsShouldStop() bool {
	switch e.dirn {
	case elevio.MD_Down:
		return (e.requests[e.floor][B_HallDown] == 1) ||
			(e.requests[e.floor][B_Cab] == 1) ||
			!e.requestBelow()

	case elevio.MD_Up:
		return (e.requests[e.floor][B_HallUp] == 1) ||
			(e.requests[e.floor][B_Cab] == 1) ||
			!e.requestAbove()

	case elevio.MD_Stop:
		fallthrough
	default:
		return true
	}
}

// request should clear immediately

// Already implemented in the function buttonPress().

// request clear at current floor
// Function implemented for CV_All
func (e *elevator) requestsClearAtCurrentFloor() {
	for btn := 0; btn < 3; btn++ {
		e.requests[e.floor][btn] = 0
	}
}

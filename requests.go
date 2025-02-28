package main

import (
	"Driver-go/elevio"
)

type Twin struct {
	m_dirn     elevio.MotorDirection
	m_behavior ElevatorBehavior
}

// Check if there are requests above the current floor
func (e Elevator) RequestsAbove() bool {
	for _floor := e.m_floor + 1; _floor < 4; _floor++ {
		for _btn := 0; _btn < 3; _btn++ {
			if e.m_requests[_floor][_btn] == 1 {
				return true
			}
		}
	}
	return false
}

// Check if there are requests below the current floor
func (e Elevator) RequestsBelow() bool {
	for _floor := 0; _floor < e.m_floor; _floor++ {
		for _btn := 0; _btn < 3; _btn++ {
			if e.m_requests[_floor][_btn] == 1 {
				return true
			}
		}
	}
	return false
}

// Check if there is a request at the current floor
func (e Elevator) RequestsHere() bool {
	for _btn := 0; _btn < 3; _btn++ {
		if e.m_requests[e.m_floor][_btn] == 1 {
			return true
		}
	}
	return false
}

// Determine the next direction based on current requests
func (e Elevator) determineDirection() Twin {
	switch e.m_dirn {
	case elevio.MD_Up:
		if e.RequestsAbove() {
			return Twin{elevio.MD_Up, EB_Moving}
		} else if e.RequestsHere() {
			return Twin{elevio.MD_Down, EB_DoorOpen}
		} else if e.RequestsBelow() {
			return Twin{elevio.MD_Down, EB_Moving}
		} else {
			return Twin{elevio.MD_Stop, EB_Idle}
		}

	case elevio.MD_Down:
		if e.RequestsBelow() {
			return Twin{elevio.MD_Down, EB_Moving}
		} else if e.RequestsHere() {
			return Twin{elevio.MD_Up, EB_DoorOpen}
		} else if e.RequestsAbove() {
			return Twin{elevio.MD_Up, EB_Moving}
		} else {
			return Twin{elevio.MD_Stop, EB_Idle}
		}

	case elevio.MD_Stop:
		if e.RequestsHere() {
			return Twin{elevio.MD_Stop, EB_DoorOpen}
		} else if e.RequestsAbove() {
			return Twin{elevio.MD_Up, EB_Moving}
		} else if e.RequestsBelow() {
			return Twin{elevio.MD_Down, EB_Moving}
		} else {
			return Twin{elevio.MD_Stop, EB_Idle}
		}
	default:
		return Twin{elevio.MD_Stop, EB_Idle} // Must include default to avoid missing return error
	}
}

// Determine if the elevator should stop at the current floor
func (e Elevator) shouldStopAtCurrentFloor() bool {
	switch e.m_dirn {
	case elevio.MD_Down:
		return (e.m_requests[e.m_floor][B_HallDown] == 1) ||
			(e.m_requests[e.m_floor][B_Cab] == 1) ||
			!e.RequestsBelow()

	case elevio.MD_Up:
		return (e.m_requests[e.m_floor][B_HallUp] == 1) ||
			(e.m_requests[e.m_floor][B_Cab] == 1) ||
			!e.RequestsAbove()

	case elevio.MD_Stop:
		fallthrough
	default:
		return true
	}
}

// Clear requests at the current floor
func (e *Elevator) clearRequestsAtCurrentFloor() {
	for _btn := 0; _btn < 3; _btn++ {
		e.m_requests[e.m_floor][_btn] = 0
	}
}

package main

import (
	"Driver-go/elevio"
	"time"
)

const (
	CV_All    = 0
	CV_InDirn = 1
)

type Elevator struct {
	floor    int
	dirn     elevio.MotorDirection
	requests [][]bool
	config   struct {
		clearRequestVariant int
	}
	behaviour ElevatorBehaviour
}

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

func main() {
	numFloors := 4
	elevio.Init("localhost:15657", numFloors)

	e := Elevator{
		floor:    0,
		dirn:     elevio.MD_Stop,
		requests: make([][]bool, numFloors),
	}
	for i := range e.requests {
		e.requests[i] = make([]bool, 3) // 3 button types
	}
	e.config.clearRequestVariant = CV_InDirn

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	doorTimer := time.NewTimer(0)
	<-doorTimer.C // Initialize timer

	for {
		select {
		case btn := <-drv_buttons:
			e.requests[btn.Floor][btn.Button] = true
			elevio.SetButtonLamp(btn.Button, btn.Floor, true)

			if e.behaviour == EB_Idle {
				dirnBehaviour := chooseDirection(e)
				e.dirn = dirnBehaviour.dirn
				e.behaviour = dirnBehaviour.behaviour

				if e.behaviour == EB_Moving {
					elevio.SetMotorDirection(e.dirn)
				} else if e.behaviour == EB_DoorOpen {
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(3 * time.Second)
				}
			}

		case floor := <-drv_floors:
			e.floor = floor
			elevio.SetFloorIndicator(floor)

			if shouldStop(e) {
				e = clearAtCurrentFloor(e)
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
				e.behaviour = EB_DoorOpen
				doorTimer.Reset(3 * time.Second)
			}

		case <-doorTimer.C:
			if e.behaviour == EB_DoorOpen {
				elevio.SetDoorOpenLamp(false)
				dirnBehaviour := chooseDirection(e)
				e.dirn = dirnBehaviour.dirn
				e.behaviour = dirnBehaviour.behaviour

				if e.behaviour == EB_Moving {
					elevio.SetMotorDirection(e.dirn)
				} else {
					e.behaviour = EB_Idle
				}
			}

		case obstructed := <-drv_obstr:
			if obstructed && e.behaviour == EB_DoorOpen {
				doorTimer.Reset(3 * time.Second)
			}

		case <-drv_stop:
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					e.requests[f][b] = false
					elevio.SetButtonLamp(b, f, false)
				}
			}
			elevio.SetMotorDirection(elevio.MD_Stop)
			e.dirn = elevio.MD_Stop
			e.behaviour = EB_Idle
		}
	}
}

type DirnBehaviourPair struct {
	dirn      elevio.MotorDirection
	behaviour ElevatorBehaviour
}

func requestsAbove(e Elevator) bool {
	for f := e.floor + 1; f < len(e.requests); f++ {
		for btn := 0; btn < 3; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < 3; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsHere(e Elevator) bool {
	for btn := 0; btn < 3; btn++ {
		if e.requests[e.floor][btn] {
			return true
		}
	}
	return false
}

func chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case elevio.MD_Up:
		if requestsAbove(e) {
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_DoorOpen}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		}
	case elevio.MD_Down:
		if requestsBelow(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{elevio.MD_Up, EB_DoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		}
	case elevio.MD_Stop:
		if requestsHere(e) {
			return DirnBehaviourPair{elevio.MD_Stop, EB_DoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		}
	}
	return DirnBehaviourPair{elevio.MD_Stop, EB_Idle}
}

func shouldStop(e Elevator) bool {
	if e.behaviour != EB_Moving {
		return false
	}

	switch e.dirn {
	case elevio.MD_Down:
		return e.requests[e.floor][elevio.BT_HallDown] ||
			e.requests[e.floor][elevio.BT_Cab] ||
			!requestsBelow(e)
	case elevio.MD_Up:
		return e.requests[e.floor][elevio.BT_HallUp] ||
			e.requests[e.floor][elevio.BT_Cab] ||
			!requestsAbove(e)
	default:
		return true
	}
}

func clearAtCurrentFloor(e Elevator) Elevator {
	switch e.config.clearRequestVariant {
	case CV_All:
		for btn := 0; btn < 3; btn++ {
			e.requests[e.floor][btn] = false
			elevio.SetButtonLamp(elevio.ButtonType(btn), e.floor, false)
		}

	case CV_InDirn:
		e.requests[e.floor][elevio.BT_Cab] = false
		elevio.SetButtonLamp(elevio.BT_Cab, e.floor, false)

		switch e.dirn {
		case elevio.MD_Up:
			e.requests[e.floor][elevio.BT_HallUp] = false
			elevio.SetButtonLamp(elevio.BT_HallUp, e.floor, false)
			if !requestsAbove(e) && !e.requests[e.floor][elevio.BT_HallUp] {
				e.requests[e.floor][elevio.BT_HallDown] = false
				elevio.SetButtonLamp(elevio.BT_HallDown, e.floor, false)
			}

		case elevio.MD_Down:
			e.requests[e.floor][elevio.BT_HallDown] = false
			elevio.SetButtonLamp(elevio.BT_HallDown, e.floor, false)
			if !requestsBelow(e) && !e.requests[e.floor][elevio.BT_HallDown] {
				e.requests[e.floor][elevio.BT_HallUp] = false
				elevio.SetButtonLamp(elevio.BT_HallUp, e.floor, false)
			}

		default:
			e.requests[e.floor][elevio.BT_HallUp] = false
			e.requests[e.floor][elevio.BT_HallDown] = false
			elevio.SetButtonLamp(elevio.BT_HallUp, e.floor, false)
			elevio.SetButtonLamp(elevio.BT_HallDown, e.floor, false)
		}
	}
	return e
}

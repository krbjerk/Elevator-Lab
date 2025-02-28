package main

import (
	"Driver-go/elevio"
	"fmt"
)

type Elevator struct {
	m_floor    int
	m_dirn     elevio.MotorDirection
	m_requests [4][3]int
	m_behavior ElevatorBehavior
}

type ElevatorBehavior int

const (
	EB_Idle     ElevatorBehavior = 0
	EB_DoorOpen                  = 1
	EB_Moving                    = 2
)

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

// Global elevator instance
var g_elevator Elevator

// Initialize the elevator
func initElevator() {
	elevio.SetMotorDirection(elevio.MD_Down)
	g_elevator.m_dirn = elevio.MD_Down
	g_elevator.m_behavior = EB_Moving
}

// Handle a button press
func handleButtonPress(_btnFloor int, _btnType elevio.ButtonType) {
	fmt.Println("Button press")

	switch g_elevator.m_behavior {
	case EB_DoorOpen:
		fmt.Println("Door is open.")
		if g_elevator.m_floor == _btnFloor {
			g_timer.startTimer(3)
			fmt.Println("door timeout 1")
		} else {
			g_elevator.m_requests[_btnFloor][_btnType] = 1
			if checkTimerExpired(g_timer) {
				g_elevator.processRequest()
				fmt.Println("Acted on request.")
			}
		}

	case EB_Moving:
		g_elevator.m_requests[_btnFloor][_btnType] = 1
	case EB_Idle:
		g_elevator.m_requests[_btnFloor][_btnType] = 1
		if checkTimerExpired(g_timer) {
			g_elevator.processRequest()
			fmt.Println("Acted on request.")
		}
	}
	updateLights(g_elevator)
	printElevatorState(g_elevator)
}

// Handle elevator arriving at a floor
func handleFloorArrival(_newFloor int) {
	fmt.Println("Arrived at floor:", _newFloor)
	g_elevator.m_floor = _newFloor
	elevio.SetFloorIndicator(g_elevator.m_floor)

	if g_elevator.m_behavior == EB_Moving && g_elevator.shouldStopAtCurrentFloor() {
		fmt.Println("Stopping elevator at floor:", _newFloor)
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		g_elevator.clearRequestsAtCurrentFloor()
		g_timer.startTimer(3)
		updateLights(g_elevator)
		g_elevator.m_behavior = EB_DoorOpen
		g_elevator.m_dirn = elevio.MD_Stop
	}
	printElevatorState(g_elevator)
}

// Handle door timeout event
func handleDoorTimeout() {
	fmt.Println("Door timeout, checking requests.")

	if g_elevator.m_behavior == EB_DoorOpen {
		twin := g_elevator.determineDirection()
		g_elevator.m_dirn = twin.m_dirn
		g_elevator.m_behavior = twin.m_behavior

		switch g_elevator.m_behavior {
		case EB_DoorOpen:
			g_timer.startTimer(3)
			g_elevator.clearRequestsAtCurrentFloor()
			updateLights(g_elevator)
		case EB_Moving:
			elevio.SetMotorDirection(g_elevator.m_dirn)
			elevio.SetDoorOpenLamp(false)
		case EB_Idle:
			elevio.SetDoorOpenLamp(false)
		}
	}
	printElevatorState(g_elevator)
}

// Process elevator request
func (e *Elevator) processRequest() {
	twin := e.determineDirection()
	e.m_dirn = twin.m_dirn
	e.m_behavior = twin.m_behavior

	switch twin.m_behavior {
	case EB_DoorOpen:
		e.clearRequestsAtCurrentFloor()
		updateLights(g_elevator)
	case EB_Moving:
		elevio.SetMotorDirection(e.m_dirn)
		elevio.SetDoorOpenLamp(false)
	case EB_Idle:
		elevio.SetDoorOpenLamp(false)
	}
}

// Update elevator lights
func updateLights(_e Elevator) {
	var BTNS = []elevio.ButtonType{elevio.BT_HallUp, elevio.BT_HallDown, elevio.BT_Cab}
	for _floor := 0; _floor < 4; _floor++ {
		for _, _btn := range BTNS {
			elevio.SetButtonLamp(_btn, _floor, convertIntToBool(_e.m_requests[_floor][_btn]))
		}
	}
}

// Convert integer to boolean
func convertIntToBool(_i int) bool {
	return _i != 0
}

// Convert direction to string
func directionToString(_dirn elevio.MotorDirection) string {
	switch _dirn {
	case 1:
		return "Up"
	case -1:
		return "Down"
	case 0:
		return "Stop"
	default:
		return "Unknown"
	}
}

// Convert behavior to string
func behaviorToString(_behavior ElevatorBehavior) string {
	switch _behavior {
	case EB_Idle:
		return "Idle"
	case EB_DoorOpen:
		return "DoorOpen"
	case EB_Moving:
		return "Moving"
	default:
		return "Unknown"
	}
}

// Print elevator state
func printElevatorState(_e Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  | Floor = %-2d         |\n"+
			"  | Dirn  = %-10s |\n"+
			"  | Behav = %-10s |\n",
		_e.m_floor,
		directionToString(_e.m_dirn),
		behaviorToString(_e.m_behavior),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")

	for _floor := 4 - 1; _floor >= 0; _floor-- {
		fmt.Printf("  | %d", _floor)
		for _btn := 0; _btn < 3; _btn++ {
			if (_floor == 4-1 && _btn == int(B_HallUp)) || (_floor == 0 && _btn == B_HallDown) {
				fmt.Print("|     ")
			} else {
				if _e.m_requests[_floor][_btn] == 1 {
					fmt.Print("|  #  ")
				} else {
					fmt.Print("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

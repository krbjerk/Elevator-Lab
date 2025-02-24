package main

import (
	"Driver-go/elevio"
	"fmt"
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

func ButtonPress(btnFloor int, btnType elevio.ButtonType) {
	fmt.Println("Button press")
	// Requests can happen at any point.

	switch Elevator.behavior {
	case EB_DoorOpen:
		fmt.Println("Door is open.")
		// Clear if request is for same floor
		// Otherwise request will be approved and stored for later
		if Elevator.floor == btnFloor {
			// Start timer for door open
			timer.Start(3)
			fmt.Println("door timeout 1")

		} else {
			Elevator.requests[btnFloor][btnType] = 1
			Elevator.actOnRequest()

		}

	case EB_Moving:
		// Take request and store it for later
		Elevator.requests[btnFloor][btnType] = 1
	case EB_Idle:
		// Take request and act on it immediately
		Elevator.requests[btnFloor][btnType] = 1
		Elevator.actOnRequest()
		fmt.Println("Acted on request.")
		// SetAlllights(Elevator)
		// Print out something about the elevator having a new state.

	}
	setAllLights(Elevator)
	elevatorPrint(Elevator)
}

func FloorArrival(newFloor int) {

	fmt.Println("Arrived at new floor.")
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
			fmt.Println("The elevator should stop here.")
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			//elevio.SetButtonLamp() TODO: Reset the button lamp on simulator
			Elevator.requestsClearAtCurrentFloor()
			// Start open door timer
			fmt.Println("door timeout 2")
			timer.Start(3)
			// Set all lights
			setAllLights(Elevator)
			Elevator.behavior = EB_DoorOpen
			Elevator.dirn = elevio.MD_Stop
		}
	default:

	}
	// Print the new state of elevator.
	elevatorPrint(Elevator)
}

func doorTimeout() {
	fmt.Println("Door timeout. Continue on.")
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
			fmt.Println("door timeout 3")
			timer.Start(3)
			Elevator.requestsClearAtCurrentFloor()
			// Set all lights
			setAllLights(Elevator)
		case EB_Moving:
			elevio.SetMotorDirection(Elevator.dirn)
			elevio.SetDoorOpenLamp(false)

		case EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(Elevator.dirn)
		}

	default:

	}
	// Print elevator state
	elevatorPrint(Elevator)
}

func elevioDirnToString(dirn elevio.MotorDirection) string {
	switch dirn {
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

func ebToString(behaviour ElevatorBehavior) string {
	switch behaviour {
	case 0:
		return "Idle"
	case 1:
		return "DoorOpen"
	case 2:
		return "Moving"
	default:
		return "Unknown"
	}
}

// Elevator print function (Converted from C)
func elevatorPrint(es elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		es.floor,
		elevioDirnToString(es.dirn),
		ebToString(es.behavior),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")

	// Print button requests
	for f := 4 - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < 3; btn++ {
			// Skip invalid buttons
			if (f == 4-1 && btn == int(B_HallUp)) || (f == 0 && btn == B_HallDown) {
				fmt.Print("|     ")
			} else {
				if es.requests[f][btn] == 1 {
					fmt.Print("|  #  ") // Requested
				} else {
					fmt.Print("|  -  ") // Not requested
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

func (e *elevator) actOnRequest() {
	twin := e.requestDirection()
	e.dirn = twin.dirn
	e.behavior = twin.behavior
	switch twin.behavior {
	case EB_DoorOpen:
		// Do something about the doorlight
		// Start the door timer
		e.requestsClearAtCurrentFloor()
		setAllLights(Elevator)

	case EB_Moving:
		//
		elevio.SetMotorDirection(e.dirn)
		elevio.SetDoorOpenLamp(false)

	case EB_Idle:
		elevio.SetDoorOpenLamp(false)
	}
}

// We still do not activate the correct lighting.
// setAllLights function to update button lights
func setAllLights(es elevator) {
	var BTNS = []elevio.ButtonType{elevio.BT_HallUp, elevio.BT_HallDown, elevio.BT_Cab}
	for floor := 0; floor < 4; floor++ {
		i := 0
		for _, btn := range BTNS {
			i++
			elevio.SetButtonLamp(btn, floor, intToBool(es.requests[floor][btn]))
		}
	}
}

func intToBool(i int) bool {
	return i != 0 // Returns true if i is nonzero, false if i is 0
}

// The elevator runs even tho the doors are open.

// We should add some initialization. The elevator behaves weirdly when it isnt initilized.

// Clean the code

// Add the right print outs

package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	if elevio.GetFloor() == -1 {
		// Make the elevator move to an actual floor on startup. Necessary for the state machine.
		initElevator()
	}

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// Create a ticker that triggers every 500ms to check the timer
	timeoutTicker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)

			handleButtonPress(a.Floor, a.Button)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			handleFloorArrival(a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				// Make it start again?
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}

		case <-timeoutTicker.C:
			if g_timer.timedOut() {
				fmt.Println("Timed out in main.")
				handleDoorTimeout()
			}
		}
	}
}

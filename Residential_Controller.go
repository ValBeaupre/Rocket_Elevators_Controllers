package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

var floorName = []string{"BSMT 2", "BSMT 1", "Lobby", "2nd", "3rd", "4th", "5th", "6th", "7th", "8th"}
var directionNameList = []string{"GoingNowhere", "Down", "Up"}
var inOutNameList = []string{"GoingIn", "GoingOut"}
var doorStatusNameList = []string{"Idle", "Stopping", "Stopped", "Moving"}
var elevatorStatusNameList = []string{"Closed", "Closing", "Opening", "Opened"}

// Constants Definitions

var appTimeout int64 = 60000

var nbFloors = 10
var lobby = 3
var timePerFloor int64 = 1000
var delayElevatorStopping int64 = 2000
var delayDoorOpening int64 = 1000
var delayBeforeCloseDoor int64 = 5000
var delayForceClose int64 = 50000
var timeoutDoorOpen int64 = 15000
var delayMaxIdleTime int64 = 15000

//Elevator directions
var goingNowhere = 0
var down = 1
var up = 2

//People entering or leaving the elevator
var goingIn = 0
var goingOut = 1

//Elevator status
var idle = 0
var stopping = 1
var stopped = 2
var moving = 3

//Door status
var closed = 0
var closing = 1
var opening = 2
var opened = 3

//Buttons function
var addFloor = 0
var callElevator = 1
var openDoor = 2
var closeDoor = 3

//Buttons status
var inactive = 0
var active = 1

func getTimeInMilli() int64 {
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	return umillisec
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

// Counter :
type Counter struct {
	count int
}

func (c *Counter) currentValue() int {
	return c.count
}

func (c *Counter) increment() {
	c.count++
}

func insertInList(list []int, insertPos int, value int) []int {
	var listTemp []int
	for index := len(list) - 1; index >= insertPos-1; index-- {
		listTemp = append(listTemp, list[index])
		list = append(list[:index], list[index+1:]...)
	}
	list = append(list, value)

	for index := len(listTemp) - 1; index >= 0; index-- {
		list = append(list, listTemp[index])
	}

	return list
}

func findIndexInList(list []int, value int) int {
	for index := 0; index <= len(list); index++ {
		if list[index] == value {
			return index
		}
	}
	return -1
}

// DirectionButton :
type DirectionButton struct {
	ID        int
	name      string
	function  int
	floor     int
	state     int
	direction int
}

// FloorButton :
type FloorButton struct {
	ID       int
	name     string
	function int
	floor    int
	state    int
}

// Door :
type Door struct {
	ID     int
	alarm  bool
	status int
}

// Elevator :
type Elevator struct {
	ID                  int
	currentFloor        int
	floorLevelForTiming int
	direction           int
	status              int
	alarm               bool
	moveTimeStamp       int64
	idleTime            int64
	openDoorTime        int64
	closeDoorTime       int64
	stopTime            int64
	forceCloseTime      int64
	destinationList     []int
	directionList       []int
	inOutList           []int
	door                *Door
	Buttons             []*FloorButton
}

// Column :
type Column struct {
	ID               int
	elevators        []*Elevator
	directionButtons []*DirectionButton
}

// Elevator Controller :
type ElevatorController struct {
	ID      int
	columns *Column
}

func (ctrl *ElevatorController) requestElevator(floorNumber int, direction int) *Elevator {
	fmt.Println(join("Elevator requested at ", floorName[floorNumber-1], " floor to go ", directionNameList[direction]))
	elevator := ctrl.callElevator(direction, floorNumber, goingIn)
	return elevator
}

func (ctrl *ElevatorController) requestFloor(elevator *Elevator, requestedFloor int) {
	ctrl.addDestination(requestedFloor, elevator, goingOut, -1)
	return
}

func (ctrl *ElevatorController) callElevator(requestedDirection int, requestedFloor int, isGoingInOrOut int) *Elevator {
	elevator := ctrl.findElevator(requestedDirection, requestedFloor)
	if ctrl.checkIfExist(requestedFloor, elevator) == false {
		ctrl.addDestination(requestedFloor, elevator, isGoingInOrOut, requestedDirection)
	}
	return elevator
}

func (ctrl *ElevatorController) addDestination(floor int, elevator *Elevator, isGoingInOrOut int, requestedDirection int) {
	if len(elevator.destinationList) > 0 {
		for _, destination := range elevator.destinationList {
			if elevator.direction != 0 {
				if elevator.direction == up {
					if (floor < destination) && (floor > elevator.currentFloor) {
						insertPosition := findIndexInList(elevator.destinationList, destination)
						elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
						elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
						elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
						return
					} else if (floor > destination) && (floor < elevator.currentFloor) {
						insertPosition := findIndexInList(elevator.destinationList, destination)
						elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
						elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
						elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
						return
					}
				} else if elevator.direction == down {
					if (floor > destination) && (floor < elevator.currentFloor) {
						insertPosition := findIndexInList(elevator.destinationList, destination)
						elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
						elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
						elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
						return
					} else if (floor < destination) && (floor > elevator.currentFloor) {
						insertPosition := findIndexInList(elevator.destinationList, destination)
						elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
						elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
						elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
						return
					}
				}
			} else {
				elevator.destinationList = append(elevator.destinationList, floor)
				elevator.directionList = append(elevator.directionList, requestedDirection)
				elevator.inOutList = append(elevator.inOutList, isGoingInOrOut)
				return
			}
		}
	} else {
		elevator.destinationList = append(elevator.destinationList, floor)
		elevator.directionList = append(elevator.directionList, requestedDirection)
		elevator.inOutList = append(elevator.inOutList, isGoingInOrOut)
		return
	}
	return
}

func (ctrl *ElevatorController) findElevator(requestedDirection int, requestedFloor int) *Elevator {
	var idleElevatorList []*Elevator
	var elevatorToUse *Elevator
	for _, elevator := range ctrl.columns.elevators {
		if elevator.door.alarm == false {
			if requestedFloor == elevator.currentFloor {
				if (elevator.status == stopped && elevator.direction == requestedDirection) || (elevator.status == idle && len(elevator.destinationList) == 0) {
					return elevator
				}
			} else if (requestedFloor > elevator.currentFloor) && ((elevator.direction == up) || (elevator.direction == goingNowhere)) && (requestedDirection == up) && ((elevator.status == moving) || (elevator.status == stopped)) {
				return elevator
			} else if (requestedFloor > elevator.currentFloor) && ((elevator.direction == down) || (elevator.direction == goingNowhere)) && (requestedDirection == down) && ((elevator.status == moving) || (elevator.status == stopped)) {
				return elevator
			} else if elevator.status == idle && len(elevator.destinationList) == 0 {
				idleElevatorList = append(idleElevatorList, elevator)
			}
		}
	}
	if len(idleElevatorList) != 0 {
		if len(idleElevatorList) > 1 {
			gap := 999999
			for _, elevator := range idleElevatorList {
				if gap > int(math.Abs(float64(elevator.currentFloor-requestedFloor))) {
					gap = int(math.Abs(float64(elevator.currentFloor - requestedFloor)))
					elevatorToUse = elevator
				}
			}
			return elevatorToUse
		} else {
			return idleElevatorList[0]
		}
	}
	elevator := ctrl.nearestElevator(requestedFloor, requestedDirection)
	if elevator != nil {
		return elevator
	} else {
		return ctrl.shortestFloorList()
	}
	return nil
}

func (ctrl *ElevatorController) verifyDestinationList() {
	for _, elevator := range ctrl.columns.elevators {
		if (len(elevator.destinationList) > 0) && (len(elevator.inOutList) > 0) {
			if elevator.destinationList[0] != 0 && elevator.status == idle {
				if elevator.currentFloor != elevator.destinationList[0] && elevator.door.status == closed {
					elevator.startMove()
				} else if (elevator.inOutList[0] == goingIn) && (elevator.currentFloor == elevator.destinationList[0]) && (elevator.door.status == closed) {
					elevator.destinationList = append(elevator.destinationList[:0], elevator.destinationList[1:]...)
					elevator.directionList = append(elevator.directionList[:0], elevator.directionList[1:]...)
					elevator.inOutList = append(elevator.inOutList[:0], elevator.inOutList[1:]...)
					elevator.openDoor()
				}
			}
		}
	}
	return
}

func (ctrl *ElevatorController) checkElevatorStatus() {
	for _, elevator := range ctrl.columns.elevators {
		if (getTimeInMilli() > (elevator.openDoorTime + timeoutDoorOpen)) && elevator.door.status != closed {
			elevator.door.alarm = true
			elevator.forceCloseDoor()
		}
		if (getTimeInMilli() > (elevator.openDoorTime + delayDoorOpening)) && (elevator.door.status == opening) {
			elevator.door.status = opened
			fmt.Println(join("Elevator ", strconv.Itoa(elevator.ID), " door is opened"))
		}
		if ((getTimeInMilli() > (elevator.openDoorTime + delayBeforeCloseDoor)) && elevator.door.status == opened) || elevator.door.alarm == true {
			elevator.closeDoor()
		}
		if (getTimeInMilli() > (elevator.closeDoorTime + delayDoorOpening)) && elevator.door.status == closing {
			if elevator.door.alarm == true && elevator.door.status != opening {
				elevator.openDoor()
				fmt.Println(join("Elevator ", strconv.Itoa(elevator.ID), " door is opening again because of an obstruction (Force Close Door"))
			}
			elevator.door.status = closed
			elevator.status = idle
			elevator.idleTime = getTimeInMilli()
			fmt.Println(join("Elevator ", strconv.Itoa(elevator.ID), " door is closed 1"))
		}
		if (getTimeInMilli() > (elevator.forceCloseTime + delayForceClose)) && elevator.door.status == closing && elevator.door.alarm == false {
			elevator.status = idle
			elevator.idleTime = getTimeInMilli()
			elevator.alarm = false
			elevator.door.alarm = false
			elevator.door.status = closed
			fmt.Println(join("Elevator ", strconv.Itoa(elevator.ID), " door is closed 2"))
		}
		if (getTimeInMilli() > (elevator.stopTime + delayElevatorStopping)) && elevator.status == stopping {
			elevator.status = stopped
			elevator.direction = goingNowhere
			fmt.Println(join("Elevator ", strconv.Itoa(elevator.ID), " is stopped at ", floorName[elevator.currentFloor-1], " floor for people ", inOutNameList[elevator.inOutList[0]]))
			elevator.directionList = append(elevator.directionList[:0], elevator.directionList[1:]...)
			elevator.inOutList = append(elevator.inOutList[:0], elevator.inOutList[1:]...)
		}
	}
	return
}

func (ctrl *ElevatorController) checkMovingElevator() {
	for _, elevator := range ctrl.columns.elevators {
		if elevator.status == moving && elevator.door.alarm == false {
			if getTimeInMilli() > (elevator.moveTimeStamp + timePerFloor) {
				if elevator.currentFloor < elevator.destinationList[0] {
					elevator.floorLevelForTiming++
				} else {
					elevator.floorLevelForTiming--
				}
				elevator.moveTimeStamp = getTimeInMilli()
				if elevator.currentFloor != elevator.floorLevelForTiming {
					elevator.currentFloor = elevator.floorLevelForTiming
					fmt.Println(join("Elevator ", strconv.Itoa(elevator.ID), " is at ", floorName[elevator.currentFloor-1], " floor "))
				}
				if elevator.currentFloor == elevator.destinationList[0] {
					elevator.stopElevator()
					ctrl.clearButtons(elevator)
					elevator.destinationList = append(elevator.destinationList[:0], elevator.destinationList[1:]...)
				}
			}
		}
		if elevator.status == stopped && elevator.door.status == closed {
			elevator.openDoor()
		}
	}
	return
}

func (ctrl *ElevatorController) clearButtons(elevator *Elevator) {
	if elevator.directionList[0] > 0 {
		button := ctrl.columns.findResquestButton(elevator.directionList[0], elevator.destinationList[0])
		button.state = inactive
		fmt.Println(join("Request direction button ", floorName[elevator.currentFloor-1], " floor is inactive"))
	}
	button := elevator.findFloorButton()
	button.state = inactive
	fmt.Println(join("Elevator ", strconv.Itoa(elevator.ID), " ", floorName[elevator.currentFloor-1], " floor button is inactive"))
	return
}

func (ctrl *ElevatorController) shortestFloorList() *Elevator {
	length := 999999
	var shortestList *Elevator
	for _, elevator := range ctrl.columns.elevators {
		if length > len(elevator.destinationList) {
			length = len(elevator.destinationList)
			shortestList = elevator
		}
	}
	return shortestList
}

func (ctrl *ElevatorController) nearestElevator(requestedFloor int, requestedDirection int) *Elevator {
	gap := 999999
	var shortestGap *Elevator
	for _, elevator := range ctrl.columns.elevators {
		if (gap > int(math.Abs(float64(elevator.currentFloor-elevator.destinationList[0])))) && elevator.door.alarm == false {
			if ((requestedFloor > elevator.currentFloor) && ((elevator.direction == up) || (elevator.direction == goingNowhere)) && (requestedDirection == up)) || ((requestedFloor > elevator.currentFloor) && ((elevator.direction == down) || (elevator.direction == goingNowhere)) && (requestedDirection == down)) {
				gap = int(math.Abs(float64(elevator.currentFloor - elevator.destinationList[0])))
				shortestGap = elevator
			}
		}

	}
	return shortestGap
}

func (ctrl *ElevatorController) checkIfExist(floor int, elevator *Elevator) bool {
	if len(elevator.destinationList) > 0 {
		for _, destination := range elevator.destinationList {
			if destination == floor {
				return true
			}
		}
	}
	return false
}

func (col *Column) findResquestButton(requestedDirection int, requestedFloor int) *DirectionButton {
	for _, button := range col.directionButtons {
		if button.direction == requestedDirection && button.floor == requestedFloor {
			return button
		}
	}
	return nil
}

func (elev *Elevator) startMove() {
	if elev.currentFloor < elev.destinationList[0] {
		elev.direction = up
	} else {
		elev.direction = down
	}
	elev.status = moving
	elev.floorLevelForTiming = elev.currentFloor
	fmt.Println(join("Elevator ", strconv.Itoa(elev.ID), " is moving ", directionNameList[elev.direction], " to ", floorName[elev.destinationList[0]-1], " floor"))
}

func (elev *Elevator) stopElevator() {
	elev.stopTime = getTimeInMilli()
	elev.status = stopping
	fmt.Println(join("Elevator ", strconv.Itoa(elev.ID), " is stopping at ", floorName[elev.destinationList[0]-1], " floor"))
}

func (elev *Elevator) openDoor() {
	elev.door.status = opening
	elev.openDoorTime = getTimeInMilli()
	fmt.Println(join("Elevator ", strconv.Itoa(elev.ID), " door is opening"))
}

func (elev *Elevator) closeDoor() {
	elev.door.status = closing
	elev.closeDoorTime = getTimeInMilli()
	fmt.Println(join("Elevator ", strconv.Itoa(elev.ID), " door is closing"))
}

func (elev *Elevator) forceCloseDoor() {
	elev.alarm = true
	elev.forceCloseTime = getTimeInMilli()
	fmt.Println(join("Elevator ", strconv.Itoa(elev.ID), " door is closing slowly (Force close)"))
}

func (elev *Elevator) findFloorButton() *FloorButton {
	for _, button := range elev.Buttons {
		if button.floor == elev.currentFloor {
			return button
		}
	}
	return nil
}

func createcontroller(nbElevators int, nbFloors int, defaultFloor [2]int) *ElevatorController {
	idGen := &Counter{count: 0}
	ctrl := &ElevatorController{ID: 1}
	ctrl.columns = &Column{ID: 1}
	if nbElevators > 0 && nbFloors > 0 && defaultFloor[0] > 0 && defaultFloor[1] > 0 {

		for index := 0; index < nbElevators; index++ {
			ctrl.columns.elevators = append(ctrl.columns.elevators, &Elevator{ID: index + 1, currentFloor: defaultFloor[index], floorLevelForTiming: defaultFloor[index], direction: goingNowhere, status: idle, alarm: false})
			ctrl.columns.elevators[index].idleTime = getTimeInMilli()
			for index2 := 0; index2 < nbFloors; index2++ {
				ctrl.columns.elevators[index].Buttons = append(ctrl.columns.elevators[index].Buttons, &FloorButton{ID: index2, name: floorName[index2], function: addFloor, state: inactive, floor: index2 + 1})
			}
			ctrl.columns.elevators[index].Buttons = append(ctrl.columns.elevators[index].Buttons, &FloorButton{ID: nbFloors + 1, name: "", function: openDoor, state: inactive})
			ctrl.columns.elevators[index].Buttons = append(ctrl.columns.elevators[index].Buttons, &FloorButton{ID: nbFloors + 1, name: "", function: closeDoor, state: inactive})
			ctrl.columns.elevators[index].door = &Door{ID: 1}
		}
		for index2 := 0; index2 < nbFloors; index2++ {
			for index3 := 0; index3 < 2; index3++ {
				if index3 == 0 {
					ctrl.columns.directionButtons = append(ctrl.columns.directionButtons, &DirectionButton{ID: idGen.currentValue(), name: string(index2 + 1), function: callElevator, floor: index2 + 1, state: inactive, direction: up})
				} else {
					ctrl.columns.directionButtons = append(ctrl.columns.directionButtons, &DirectionButton{ID: idGen.currentValue(), name: string(index2 + 1), function: callElevator, floor: index2 + 1, state: inactive, direction: down})
				}
			}

		}

	}
	return ctrl
}

/*
FLOOR INDEX:
    BSMT 2 : 1
    BSMT 1 : 2
    Lobby : 3
    2nd : 4
    3rd : 5
    4th : 6
    5th : 7
    6th : 8
    7th : 9
    8th : 10

ELEVATOR INDEX:
    1 : 1
    2 : 2 */

func main() {
	var defaultFloor [2]int
	var elevator [5]*Elevator
	var firstInstructionDone = false
	var timeout = getTimeInMilli()
	var readyToStop = false

	defaultFloor[0] = 10
	defaultFloor[1] = 3
	ctrl := createcontroller(2, 10, defaultFloor)

	for {
		if firstInstructionDone == false {
			firstInstructionDone = true
			elevator[0] = ctrl.requestElevator(10, down)
			elevator[1] = ctrl.requestElevator(3, down)
			//elevator[2] = ctrl.requestElevator(9, down)

			ctrl.requestFloor(elevator[0], 3)
			//ctrl.requestFloor(elevator[0], 5)
			//ctrl.requestFloor(elevator[0], 2)
			ctrl.requestFloor(elevator[1], 2)
			//ctrl.requestFloor(elevator[1], 5)
			//ctrl.requestFloor(elevator[1], 8)
			//ctrl.requestFloor(elevator[2], 2)
		}

		ctrl.checkElevatorStatus()
		ctrl.verifyDestinationList()
		ctrl.checkMovingElevator()

		if (getTimeInMilli() > (ctrl.columns.elevators[0].idleTime + appTimeout)) && (getTimeInMilli() > (ctrl.columns.elevators[1].idleTime + appTimeout)) {
			readyToStop = true
		}
		if (getTimeInMilli() > (timeout + appTimeout)) || readyToStop == true {
			//fmt.Println(string(appTimeout))
			//fmt.Println("Break")
			break
		}
	}
}

/*
*************************************************************************************
Residential Scenario 1:
A User located at Floor 1 calls for elevators Originating from Basement 1 and Floor 4, he gets one and gets to the Fifth floor with it'

Residential Scenario 2:
A User located at Floor Basement 2 calls for elevators Originating from Floor 8 and Floor 1, he gets one to get to the 4th floor,
Simultaneously, someone at Floor 1 requests an elevator to get to 3rd floor AS someone at 7th requests to go down to basement 1'

Residential Scenario 3:
A User located at Floor 8 calls for elevators Originating from Floor 8 and Floor 1, he gets one to get to the 1st floor.
Simultaneously Someone at Floor 1 requests an elevator to get to Basement 1'
*/

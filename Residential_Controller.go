package main

import (
	"math"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var floorName = []string{"BSMT 2", "BSMT 1", "Lobby", "2nd", "3rd", "4th", "5th", "6th", "7th", "8th"}
var directionNameList = []string{"GoingNowhere", "Down", "Up"}
var inOutNameList = []string{"GoingIn", "GoingOut"}
var doorStatusNameList = []string{"Idle", "Stopping", "Stopped", "Moving"}
var elevatorStatusNameList = []string{"Closed", "Closing", "Opening", "Opened"}

//Constants definition
const (
	nbFloors              int   = 10
	lobby                 int   = 3
	timePerFloor          int64 = 1000
	delayElevatorStopping int64 = 2000
	delayDoorOpening      int64 = 1000
	delayBeforeCloseDoor  int64 = 5000
	delayForceClose       int64 = 50000
	timeoutDoorOpen       int64 = 15000
	delayMaxIdleTime      int64 = 15000
	appTimeout            int64 = 60000

	//Elevator directions
	goingNowhere int = 0
	down         int = 1
	up           int = 2

	//People entering or leaving the elevator
	goingIn  int = 0
	gointOut int = 1

	//Elevator status
	idle     int = 0
	stopping int = 1
	stopped  int = 2
	moving   int = 3

	//Door status
	closed  int = 0
	closing int = 1
	opening int = 2
	opened  int = 3

	//Buttons function
	addFloor     int = 0
	callElevator int = 1
	openDoor     int = 2
	closeDoor    int = 3

	//Buttons status
	inactive int = 0
	active   int = 1
)

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

// OpenDoorButton :
type OpenDoorButton struct {
	door int
}

// CloseDoorButton :
type CloseDoorButton struct {
	door int
}

// FloorButton :
type FloorButton struct {
	floor  int
	status int
}

// DirectionButton :
type DirectionButton struct {
	floor     int
	direction int
	status    int
}

// Door :
type Door struct {
	ID          int
	status      int
	alarm       bool
	obstruction bool
}

// Elevator :
type Elevator struct {
	ID                  int
	currentFloor        int
	direction           int
	status              int
	floorLevelForTiming int
	alarm               bool
	moveTimeStamp       int64
	idleTime            int64
	openDoorTime        int64
	closeDoorTime       int64
	stopTime            int64
	forceCloseTime      int64
	destinationList     []int
	inOutList           []int
	directionList       []int
	floorButtons        []int
	door                *Door
	openDoorButton      *OpenDoorButton
	closeDoorButton     *CloseDoorButton
}

// Column :
type Column struct {
	nbElev              int
	defaultFloor        int
	nbDirectionButtons  int
	elevator            []int
	destinationList     []int
	directionButtonList []int
}

// Floor :
type Floor struct {
	ID   int
	name int
}

// ElevatorController :
type ElevatorController struct {
	ID      int
	columns *Column
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

func (elev *Elevator) findFloorButton() {
	for _, button := range elev.floorButtons {
		if button.floor == elev.currentFloor {
			return button
		}
	}
}

func (elev *Elevator) startMove() {
	if elev.currentFloor < lev.destinationList[0] {
		elev.direction = up
	} else {
		elev.direction = down
	}
	elev.status = moving
	elev.floorLevelForTiming = elev.currentFloor
	fmt.Println(join("Elevator ", strconv.Itoa(elev.ID), " is moving ", directionNameList[elev.direction], " to ", floorName[elev.destinationList[0]-1], " floor"))
}

func (col *Column) findDirectionButton() {
	for _, button := range col.directionButtonList {
		if button.direction == requestedDirection && button.floor == requestedFloor {
			return button
		}
	}
	return nil
}

func (ctrl *ElevatorController) requestElevator(floorNumber int, direction int) *Elevator {
	fmt.Println(join("Elevator requested at ", floorNames[floorNumber-1], " to go ", directionNameList[direction]))
	elevator := ctrl.callElevator(direction, floorNumber, goingIn)
	return elevator
}

func (ctrl *ElevatorController) requestFloor(elevator, requestedFloor) {
	ctrl.addDestinationElev(requestedFloor, elevator, gointOut, -1)
	return
}

func (ctrl *ElevatorController) callElevator(requestedDirection int, requestedFloor int, isGoingInOrOut int) {
	elevator := ctrl.findElevator(requestedDirection, requestedFloor)
	if ctrl.checkIfDestinExist(requestedFloor, elevator) == false {
		ctrl.addDestinationElev(requestedFloor, elevator, gointOut, requestedDirection)
	}
	return
}

func (ctrl *ElevatorController) checkIfDestinExist(floor int, elevator *Elevator) bool {
	if len(elevator.destinationList) > 0 {
		for _, destination := range elevator.destinationList {
			if destination == floor {
				return true
			}
		}
	}
	return false
}

func (ctrl *ElevatorController) shortestdestinationList() *Elevator {
	length := 99999
	var shortestList *Elevator
	for _, elevator := ctrl.columns.elevators {
		if length > len(elevator.destinationList) {
			length = len(elevator.destinationList)
			elevWithShortestList = elevator
		}
	}
	return elevWithShortestList
}

func (ctrl *ElevatorController) nearestElevator(requestedFloor int, requestedDirection int) *Elevator {
	gap := 99999 
	var shortestGap *Elevator
	for _, elevator := range ctrl.columns.elevators {
		if (gap > int(math.Abs(float64(elevator.currentFloor - elevator.destinationList[0])))) && elevator.door.alarm == false {
			if ((requestedFloor > elevator.currentFloor) && ((elevator.direction == up) || (elevator.direction == goingNowhere)) && (requestedDirection == up)) || ((requestedFloor > elevator.direction == down || (elevator.direction == goingNowhere)) && (requestedDirection == down)) {
				gap = int(math.Abs(float64(elevator.currentFloor - elevator.destinationList[0])))
				elevWithShortestGap = elevator
			}
		}
	}
	return elevWithShortestGap
}

func (ctrl *ElevatorController) clearButtons(elevator *ElevatorController) {
	if elevator.directionList[0] > 0 {
		button := ctrl.columns.findDirectionButton(elevator.directionList[0], elevator.destinationList[0])
		button.status = inactive
		fmt.Println(join("Request button direction ", floorName[elevator.currentFloor-1], " floor is inactive"))
	}
	button := elevator.findFloorButton(elevator.destinationList[0])
	button.status = inactive
	fmt.Println(join("Elevator ", strconv.Itoa(elevator.ID), " ", floorName[elevator.currentFloor-1], " floor button is inactive"))
	return
}

func (ctrl *ElevatorController) checkMovingElevator() {
	for _, elevator := range ctrl.columns.elevators {
		if elevator.status == moving && elevator.door.alarm == false {
			if (getTimeInMilli() - elevator.moveTimeStamp) >= timePerFloor {
				if elevator.currentFloor < elevator.destinationList[0] {
					elevator.floorLevelForTiming++
				} else {
					elevator.floorLevelForTiming--
				}
				elevator.moveTimeStamp = getTimeInMilli
				if elevator.currentFloor != elevator.floorLevelForTiming {
					elevator.currentFloor = elevator.floorLevelForTiming
					fmt.Println(join("Elevator ", strconv.Itoa(elevator.ID), " is at ", floorName[elevator.currentFloor-1], "floor"))
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

func (ctrl *ElevatorController) verifyDestinationList() {
	for _, elevator := range ctrl.columns.elevators {
		if (len(elevator.destinationList) > 0) && (len(elevator.inOutList) >0) {
			if elevator.destinationList[0] != 0 && elevator.status == idle {
				if elevator.currentFloor != elevator.destinationList[0] && elevator.door.status == closed {
					elevator.startMove()
				} else if (elevator.inOutList[0] && elevator.door.status == closed)
				elevator.destinationList = append(elevator.destinationList[:0], elevator.destinationList[1:]...)
				elevator.directionList = append(elevator.directionList[:0], elevator.directionList[1:]...)
				elevator.inOutList = append(elevator.inOutList[:0], elevator.inOutList[1:]...)
				elevator.openDoor()
			}
		}
	}
	return
}

func (ctrl *ElevatorController) checkElevatorStatus()

/*
def checkElevatorStatus(self):
	for elevator in self.columns.elevators:
		if (timeInMilli() > elevator.openDoorTime + timeOutDoorOpen) and (elevator.door.status is not Closed):
			elevator.door.alarm = True
			elevator.forceCloseDoor()
		if (timeInMilli() > elevator.idleTime + delayMaxIdleTime) and (not elevator.destinationList):
			pass
			#self.AddDestination(self.Columns.DefaultFloor,Elevator,GoingOut,"")

		if (timeInMilli() > (elevator.openDoorTime + delayDoorOpening)) and (elevator.door.status is Opening):
			elevator.door.status = Opened
			print("Elevator " + str(elevator.ID) + " door is opened") 
		
		if ((timeInMilli() > (elevator.openDoorTime + delayBeforeCloseDoor)) and (elevator.door.status is Opened)) or elevator.door.alarm is True:
			elevator.closeDoor()
		
		if (timeInMilli() > elevator.closeDoorTime + delayDoorOpening) and elevator.door.status is Closing:
			if (elevator.door.alarm is True) and (elevator.door.status is not Opening):
				elevator.openDoor()
				print("Elevator " + str(elevator.ID) + " door is opening again because of obstruction")
			elevator.door.status = Closed
			elevator.status = Idle
			elevator.idleTime = timeInMilli()
			print("Elevator " + str(elevator.ID) + " door is closed")
		
		if (timeInMilli() > elevator.forceCloseTime + delayForceCloseDoor) and (elevator.door.status is Closing) and elevator.door.alarm is False:
			elevator.status = Idle
			elevator.idleTime = timeInMilli()
			elevator.alarm = False
			elevator.door.alarm = False
			elevator.door.status = Closed
			print("Elevator " + str(elevator.ID) + " door is closed")

		if (timeInMilli() > elevator.stopTime + delayElevatorStopping) and elevator.status is Stopping:
			elevator.status = Stopped
			elevator.direction = GoingNowhere
			print("Elevator " + str(elevator.ID) + " is stopped at " + str(floorNames[elevator.currentFloor-1]) + " floor for people " + inOrOutNameList[elevator.inOutList[0]]) 
			elevator.directionList.pop(0)
			elevator.inOutList.pop(0)

*/

func (ctrl *ElevatorController) addDestinationElev(floor int, elevator *Elevator, isGoingInOrOut int, requestedDirection int) {
	if len(elevator.destinationList) > 0 {
		for _, destination := range elevator.destinationList {
			if elevator.direction != 0 {
				if elevator.direction == up {
					if (floor < destination) && (floor > elevator.currentFloor) {
						// insertPosition := elevator.destinationList = append(elevator.destinationList[:0], append([]T{x}, a[1:]...)...)
						insertPosition := findIndexInList(elevator.destinationList, destination)
						elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
						elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
						elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
						return
					} else if (floor > destination) && (floor < elevator.currentFloor) {
						insertPosition := findIndexInList(elevator.destinationList, destination)
						elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
						elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
						elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
						return
					}
				} else if elevator.direction == down {
					if (floor > destination) && (floor < elevator.currentFloor) {
						insertPosition := findIndexInList(elevator.destinationList, destination)
						elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
						elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
						elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
						return
					} else if (floor < destination) && (floor > elevator.currentFloor) {
						insertPosition := findIndexInList(elevator.destinationList, destination)
						elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
						elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
						elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
						return
					}
				}
			} else {
				elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
				elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
				elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
				return
			}
		}
	} else {
		elevator.destinationList = insertInList(elevator.destinationList, insertPosition, floor)
		elevator.inOutList = insertInList(elevator.inOutList, insertPosition, isGoingInOrOut)
		elevator.directionList = insertInList(elevator.directionList, insertPosition, requestedDirection)
		return
	}
	return
}



/*
*/

/*
def findElevator(self, requestedDirection, requestedFloor):
	idleElevatorList = []
	for elevator in self.columns.elevators:
		if elevator.door.alarm is False:
			if requestedFloor is elevator.currentFloor:
				if (elevator.status is Stopped and elevator.direction is requestedDirection) or (elevator.status is Idle and not elevator.destinationList):
					return elevator
			elif ((requestedFloor > elevator.currentFloor) and ((elevator.direction is Up) or (elevator.direction is GoingNowhere)) and (requestedDirection is Up) and ((elevator.status is Moving) or (elevator.status is Stopped))):
				return elevator
			elif ((requestedFloor > elevator.currentFloor) and ((elevator.direction is Down) or (elevator.direction is GoingNowhere)) and (requestedDirection is Down) and ((elevator.status is Moving) or (elevator.status is Stopped))):
				return elevator
			elif elevator.status is Idle and not elevator.destinationList:
				idleElevatorList.append(elevator)
	if idleElevatorList:
		if (len(idleElevatorList) > 1):
			gap = 99999
			for elevator in idleElevatorList:
				if gap > abs(elevator.currentFloor - requestedFloor):
					gap = abs(elevator.currentFloor - requestedFloor)
					elevatorToUse = elevator
			return elevatorToUse
		else:
			return idleElevatorList[0]

	elevator = self.nearestElevator(requestedFloor, requestedDirection)
	if not elevator:
		return elevator
	else:
		return self.shortestdestinationList()
*/

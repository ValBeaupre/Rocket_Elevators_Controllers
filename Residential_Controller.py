# Residential Elevator Controller In Python From My Pseudocode From Last Week
# Written by : Valerie Beaupre
# Date : 2018/09/24 TO 2018/09/27

import itertools
import os
import time
import datetime
from datetime import timedelta

nbColumn = 1
nbElev = 2
nbFloors = 10
floorNames = ["BSMT 1", "BSMT 2", "Lobby", "2nd","3rd", "4th", "5th", "6th", "7th", "8th"]
lobby = 3
defaultFloor = 3
timePerFloor = 1000

# Elevator Directions
directionNameList = ["GoingNowhere", "Down", "Up"]
GoingNowhere = 0
Down = 1
Up = 2

# People entering or leaving the elevator
inOrOutNameList = ["GoingIn","GoingOut"]
GoingIn = 0
GoingOut = 1

#Elevator status
elevatorStatusNameList = ["Idle","Stopping","Stopped","Moving"]
Idle = 0
Stopping = 1
Stopped = 2
Moving = 3

#Door status
oorStatusNameList= ["Closed","Closing","Opening","Opened"]
Closed = 0
Closing = 1
Opening = 2
Opened = 3

#Buttons Function
AddFloor = 0
CallElevator = 1
OpenDoor = 2
CloseDoor = 3

#Buttons status
Inactive = 0
Active = 1

delayElevatorStopping = 2000
delayDoorOpening = 1000
delayBeforeCloseDoor = 5000
delayForceCloseDoor = 50000
timeOutDoorOpen = 15000
delayMaxIdleTime = 15000   
currentFloor = 0
appTimeOut = 60000

timeInMilli = lambda: int(round(time.time() * 1000))

def counter(): #Self increment counter to generate unique ID's
    return lambda c=itertools.count(): next(c)

count = counter()  


class ElevatorController:
    def __init__(self):
        self.id = 1
        self.columns = Column(2, lobby, 2)
        print("Controller Created a column for " + str(nbElev) + " elevators and for " + str(nbFloors) + " floors")

    def requestElevator(self, floorNumber, direction):
        print("Elevator requested at " + floorNames[floorNumber-1] + " to go " + directionNameList[direction])
        elevator = self.callElevator(direction, floorNumber, GoingIn)
        return elevator

    def requestFloor(self, elevator, requestedFloor):
        self.addDestinationElev(requestedFloor, elevator, GoingOut,"")

    def callElevator(self, requestedDirection, requestedFloor, isGoingInOrOut):
        elevator = self.findElevator(requestedDirection, requestedFloor)
        if self.checkIfDestinExist(requestedFloor, elevator) is False:
            self.addDestinationElev(requestedFloor, elevator, isGoingInOrOut, requestedDirection)
        return elevator

    def checkIfDestinExist(self, floor, elevator):
        if elevator.destinationList:
            for destination in elevator.destinationList:
                if destination == floor:
                    return True
        return False

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

    def shortestdestinationList(self):
        length = 99999
        #elevWithShortestList = None
        for elevator in self.columns.elevators:
            if length > len(elevator.destinationList):
                length = len(elevator.destinationList)
                elevWithShortestList = elevator
        return elevWithShortestList
        print("shortestList")

    def nearestElevator(self, requestedFloor, requestedDirection):
        gap = 99999
        elevWithShortestGap = None
        for elevator in self.columns.elevators:
            if (gap > abs(elevator.currentFloor - elevator.destinationList[0])) and (elevator.door.alarm is False):
                if ((requestedFloor > elevator.currentFloor) and ((elevator.direction is Up) or (elevator.direction is GoingNowhere)) and (requestedDirection is Up)) or ((requestedFloor > elevator.currentFloor) and ((elevator.direction is Down) or (elevator.direction is GoingNowhere)) and (requestedDirection is Down)):
                    gap = abs(elevator.currentFloor - elevator.destinationList[0])
                    elevWithShortestGap = elevator
        return elevWithShortestGap
        print("nearestELev")

    def addDestinationElev(self, floor, elevator, isGoingInOrOut, requestedDirection):
        if elevator.destinationList:
            for destination in elevator.destinationList:
                if elevator.direction:
                    if elevator.direction is Up:
                        if (floor < destination) and (floor > elevator.currentFloor):
                            insertPosition = elevator.destinationList.index(destination)
                            elevator.destinationList.insert(insertPosition, floor)
                            elevator.inOutList.insert(insertPosition, isGoingInOrOut)
                            elevator.directionList.insert(insertPosition, requestedDirection)
                            return
                        elif (floor > destination) and (floor < elevator.currentFloor):
                            insertPosition = elevator.destinationList.index(destination)
                            elevator.destinationList.insert(insertPosition, floor)
                            elevator.inOutList.insert(insertPosition, isGoingInOrOut)
                            elevator.directionList.insert(insertPosition, requestedDirection)
                            return
                    elif elevator.direction is Down:
                        if (floor > destination) and (floor < elevator.currentFloor):
                            insertPosition = elevator.destinationList.index(destination)
                            elevator.destinationList.insert(insertPosition, floor)
                            elevator.inOutList.insert(insertPosition, isGoingInOrOut)
                            elevator.directionList.insert(insertPosition, requestedDirection)
                            return
                        elif (floor < destination) and (floor > elevator.currentFloor):
                            insertPosition = elevator.destinationList.index(destination)
                            elevator.destinationList.insert(insertPosition, floor)
                            elevator.inOutList.insert(insertPosition, isGoingInOrOut)
                            elevator.directionList.insert(insertPosition, requestedDirection)
                            return

                else:
                    elevator.destinationList.append(floor)
                    elevator.inOutList.append(isGoingInOrOut)
                    elevator.directionList.append(requestedDirection)
                    return
        else:
            elevator.destinationList.append(floor)
            elevator.inOutList.append(isGoingInOrOut)
            elevator.directionList.append(requestedDirection)
            return

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
                print("Elevator " + str(elevator.id) + " door is opened") 
            
            if ((timeInMilli() > (elevator.openDoorTime + delayBeforeCloseDoor)) and (elevator.door.status is Opened)) or elevator.door.alarm is True:
                elevator.closeDoor()
            
            if (timeInMilli() > elevator.closeDoorTime + delayDoorOpening) and elevator.door.status is Closing:
                if (elevator.door.alarm is True) and (elevator.door.status is not Opening):
                    elevator.openDoor()
                    print("Elevator " + str(elevator.id) + " door is opening again because of obstruction")
                elevator.door.status = Closed
                elevator.status = Idle
                elevator.idleTime = timeInMilli()
                #print("CloseDoor 1")
                print("Elevator " + str(elevator.id) + " door is closed")
            
            if (timeInMilli() > elevator.forceCloseTime + delayForceCloseDoor) and (elevator.door.status is Closing) and elevator.door.alarm is False:
                elevator.status = Idle
                elevator.idleTime = timeInMilli()
                elevator.alarm = False
                elevator.door.alarm = False
                elevator.door.status = Closed
                #print("CloseDoor 2")
                print("Elevator " + str(elevator.id) + " door is closed")

            if (timeInMilli() > elevator.stopTime + delayElevatorStopping) and elevator.status is Stopping:
                elevator.Status = Stopped
                elevator.Direction = GoingNowhere
                print("Elevator " + str(elevator.id) + " is stopped at " + str(floorNames[elevator.currentFloor-1]) + " floor for people " + inOrOutNameList[elevator.inOutList[0]]) 
                elevator.directionList.pop(0)
                elevator.inOutList.pop(0)

    def verifyDestinationList(self):
        for elevator in self.columns.elevators:
            if elevator.destinationList:
                if elevator.inOutList:
                    if (elevator.destinationList[0]) and (elevator.status is Idle):
                        if (elevator.currentFloor is not elevator.destinationList[0]) and elevator.door.status is Closed:
                            elevator.startMove()
                        elif (elevator.inOutList[0] is GoingIn) and (elevator.currentFloor is elevator.destinationList[0]) and elevator.door.status is Closed:
                            print("Open Door 1")
                            elevator.destinationList.pop(0)
                            elevator.directionList.pop(0)
                            elevator.inOutList.pop(0)
                            elevator.openDoor()

    def checkMovingElevator(self):
        for elevator in self.columns.elevators:
            if (elevator.status is Moving) and (elevator.door.alarm is False):
                if (timeInMilli() - elevator.moveTimeStamp) >= timePerFloor:
                    if elevator.currentFloor < elevator.destinationList[0]:
                        elevator.floorLevelForTiming += 1
                    else:
                        elevator.floorLevelForTiming -= 1

                    elevator.moveTimeStamp = timeInMilli()
                    if elevator.currentFloor is not elevator.floorLevelForTiming:
                        elevator.currentFloor = elevator.floorLevelForTiming
                        print("Elevator " + str(elevator.id) + " is at " + floorNames[elevator.currentFloor-1] + " floor")
                        
                    if elevator.currentFloor is elevator.destinationList[0]:
                        elevator.stopElevator()
                        self.clearButtons(elevator)
                        elevator.destinationList.pop(0)
                        #print(elevator.destinationList)

            if elevator.status is Stopped and elevator.door.status is Closed:
                #print("OpenDoor 2")
                elevator.openDoor()

    def clearButtons(self, elevator):
        if elevator.directionList[0] is not "":
            button = self.columns.findDirectionButton(elevator.directionList[0],elevator.destinationList[0]) 
            button.status = Inactive
            print("Request button direction " + directionNameList[button.direction] + " at " + floorNames[elevator.destinationList[0]-1] + " floor is inactive")
        button = elevator.findFloorButton()
        button.status = Inactive
        print("Elevator " + str(elevator.id) + " " + floorNames[elevator.currentFloor-1] + " floor button is inactive")


class Floor:
    def __init__(self, id, name, buttons):
        self.id = id
        self.name = name
        #self.directionButtons = buttons


class Column:
    def __init__(self, nbElev, defaultFloor, nbDirectionButtons):
        self.nbElev = nbElev
        self.defaultFloor = defaultFloor
        self.nbDirectionButtons = nbDirectionButtons
        self.elevators = []
        for index in range(self.nbElev):
            self.elevators.append(Elevator(index+1))
        #for elevator in range(self.elevators):
            #print("elevator.id")

        self.destinationList = []
        self.directionButtons = []
        for index in range(nbFloors):
            directionButtonList = []
            if index == 0:
                directionButtonList.append(DirectionButton(Up))
            elif index == (nbFloors - 1):
                directionButtonList.append(DirectionButton(Down))
            else:
                directionButtonList.append(DirectionButton(Up))
                directionButtonList.append(DirectionButton(Down))

            self.destinationList.append(Floor(index, floorNames[index], directionButtonList))

    def findDirectionButton(self, direction, requestedFloor):
        for button in self.directionButtons:
            if (requestedFloor == self.directionButtons) and (direction == self.directionButtons):
                return button
        print("Column Created")


class Elevator:
    def __init__(self, id):
        self.id = id
        self.currentFloor = lobby
        self.direction = GoingNowhere
        self.status = Idle
        self.floorLevelForTiming = lobby
        self.alarm = False
        self.moveTimeStamp = 0
        self.idleTime = timeInMilli()
        self.openDoorTime = timeInMilli()
        self.closeDoorTime = timeInMilli()
        self.stopTime = timeInMilli()
        self.forceCloseTime = timeInMilli()
        self.destinationList = []     
        self.inOutList = []
        self.directionList = []
        self.floorButtons = []  #Buttons inside the elevators
        for index in range(nbFloors):
            self.floorButtons.append(FloorButton(index))

        self.door = Door()
        self.openDoorButton = OpenDoorButton(self.door)
        self.closeDoorButton = CloseDoorButton(self.door)
        print("Elevator Created")

    def findFloorButton (self, currentFloor):
        for button in self.floorButtons: 
            if self.currentFloor is self.floorButtons:
                return button 

    def startMove(self):
        if (self.currentFloor < self.destinationList[0]): 
            self.direction : Up
        elif (self.currentFloor > self.destinationList[0]): 
            self.direction : Down
        self.status : Moving
        self.moveTimeStamp = timeInMilli()
        self.floorLevelForTiming = self.currentFloor
        #print("Elevator " + str(self.id) + " is moving " + directionNameList[self.direction] + " to " + floorNames[self.destination[0]-1] + " floor for people " + inOrOutNameList[self.inOutList[0]])
        print("Elevator " + str(self.id) + " is moving " + directionNameList[self.direction] + " to " + floorNames[self.destinationList[0]-1] + " floor")

    def stopElevator(self):
        self.stopTime = timeInMilli()
        self.Status = Stopping
        print("Elevator " + str(self.id) + " is stopping at " + floorNames[self.destinationList[0]-1] + " floor")

    def openDoor(self):
        self.door.status = Opening
        self.OpenDoorTime = timeInMilli()
        print("Elevator " + str(self.id) + " door is opening")

    def closeDoor(self):
        self.door.status = Closing
        self.closeDoorTime = timeInMilli()
        print("Elevator " + str(self.id) + " door is closing")

    def forceCloseDoor(self):
        self.alarm = True 
        print("Elevator " + str(self.id) + " door is losing slowly (Force Close)")
        self.forceCloseTime = timeInMilli()


class DirectionButton:
    def __init__(self, direction):
        self.direction = ""  # UP or DOWN
        self.status = Inactive  # Active or Inactive


class FloorButton:
    def __init__(self, floor):
        self.floor = floor
        self.status = Inactive


class Door:
    def __init__(self):
        self.id = id
        self.status = Closed
        self.obstruction = False
        self.alarm = False


class OpenDoorButton:
    def __init__(self, door):
        self.door = door


class CloseDoorButton:
    def __init__(self, door):
        self.door = door


controller = ElevatorController()

def initTest(elev1Floor,elev2Floor):
    controller.columns.elevators[0].currentFloor = elev1Floor
    controller.columns.elevators[0].floorLevelForTiming = elev1Floor
    print("Elevator 1 current floor is " + floorNames[controller.columns.elevators[0].currentFloor-1])
    
    controller.columns.elevators[1].currentFloor = elev2Floor
    controller.columns.elevators[1].floorLevelForTiming = elev2Floor
    print("Elevator 2 current floor is " + floorNames[controller.columns.elevators[1].currentFloor-1])


"""
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
    2 : 2

STATUS INDEX:
    GOINGNOWHERE : 0
    DOWN : 1
    UP : 2

"""

#CONTROLLER INITIALISATION
initTest(10,3)
controllerRunning = True

def main():
    elevator = []
    firstInstructionDone = False
    timeOut = timeInMilli()
    readyToStopController = False 

    while controllerRunning is True:
        #INITIATE ALL THE CALLS
        if firstInstructionDone is False:
            firstInstructionDone = True
            #Firt user
            elevator.append(controller.requestElevator(10, Down))
            #Second user
            elevator.append(controller.requestElevator(1, Up))
            #Third user
            #Elevator.append(controller.requestElevator(9, Down))
            #First user
            controller.requestFloor(elevator[0],3)
            #Second user
            controller.requestFloor(elevator[1],5)
            #Third user
            #controller.requestFloor(elevator[2],1)
    

        #CONTROLLER SEQUENCES
        controller.checkElevatorStatus()
        controller.verifyDestinationList()
        controller.checkMovingElevator()


        #STOP THE CONTROLLER WHEN DONE OR TIMED OUT
        if ((timeInMilli() > controller.columns.elevators[0].idleTime + appTimeOut) and (timeInMilli() > controller.columns.elevators[1].idleTime + appTimeOut)):
            readyToStopController = True
        if (timeInMilli() > timeOut + appTimeOut) or (readyToStopController is True):
            print("Break")
            break

main()


# **********************************************************************************
# **********************************************************************************

""" ****************************************************************************************

Residential Scenario 1:
A User located at Floor 1 calls for elevators Originating from Basement 1 and Floor 4, he gets one and gets to the Fifth floor with it'

Residential Scenario 2:
A User located at Floor Basement 2 calls for elevators Originating from Floor 8 and Floor 1, he gets one to get to the 4th floor, Simultaneously, someone at Floor 1 requests an elevator to get to 3rd floor as someone at 7th requests to go down to basement 1'

Residential Scenario 3:
A User located at Floor 8 calls for elevators Originating from Floor 8 and Floor 1, he gets one to get to the 1st floor. Simultaneously Someone at Floor 1 requests an elevator to get to Basement 1'

 """

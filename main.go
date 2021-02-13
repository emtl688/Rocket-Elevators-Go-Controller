package main

import (
	"fmt"
	"math"
	"sort"
)

//--------------------GLOBAL VARIABLES--------------------//
var columnID int = 1
var floorRequestButtonID = 1
var callButtonID = 1
var buttonFloor = 1

//--------------------BATTERY CLASS--------------------//

type Battery struct {
	ID                        int
	status                    string
	amountOfColumns           int
	amountOfFloors            int
	amountOfBasements         int
	amountOfElevatorPerColumn int
	servedFloors              []int
	columnsList               []Column
	floorRequestButtonsList   []FloorRequestButton
}

func (battery *Battery) createBasementColumn(amountOfBasements int, amountOfElevatorPerColumn int) {
	var servedFloors []int
	floor := -1
	for i := 0; i < (amountOfBasements + 1); i++ {
		if i == 0 {
			servedFloors = append(servedFloors, 1)
		} else {
			servedFloors = append(servedFloors, floor)
			floor--
		}
	}
	column := Column{
		ID:                columnID,
		status:            "online",
		amountOfElevators: amountOfElevatorPerColumn,
		servedFloorsList:  servedFloors,
		isBasement:        true,
		elevatorsList:     []Elevator{},
		callButtonsList:   []CallButton{},
	}
	column.createElevators(servedFloors, column.amountOfElevators)
	column.createCallButtons(column.servedFloorsList, battery.amountOfBasements, column.isBasement)
	battery.columnsList = append(battery.columnsList, column)
	columnID++
}

func (battery *Battery) createColumns(amountOfColumns int, amountOfFloors int, amountOfBasements int, amountOfElevatorPerColumn int) {
	//We get the average number of floors per column
	amountOfFloorsPerColumn := int(math.Ceil(float64(amountOfFloors / amountOfColumns)))
	floor := 1
	for i := 1; i <= amountOfColumns; i++ {
		var servedFloors []int
		for i := 1; i <= amountOfFloorsPerColumn; i++ {
			if columnID > 2 && i == 1 {
				servedFloors = append(servedFloors, floor-1)
			}
			if floor <= amountOfFloors {
				servedFloors = append(servedFloors, floor)
				floor++
			}
		}
		column := Column{
			ID:                columnID,
			status:            "online",
			servedFloorsList:  servedFloors,
			amountOfElevators: amountOfElevatorPerColumn,
			isBasement:        false,
			elevatorsList:     []Elevator{},
			callButtonsList:   []CallButton{},
		}
		column.createElevators(servedFloors, amountOfElevatorPerColumn)
		column.createCallButtons(column.servedFloorsList, amountOfBasements, column.isBasement)
		battery.columnsList = append(battery.columnsList, column)
		columnID++
	}
}

func (battery *Battery) createFloorRequestButtons(amountOfFloors int) {
	buttonFloor := 1
	for i := 0; i < amountOfFloors; i++ {
		floorRequestButton := FloorRequestButton{
			ID:     floorRequestButtonID,
			status: "off",
			floor:  buttonFloor,
		}
		battery.floorRequestButtonsList = append(battery.floorRequestButtonsList, floorRequestButton)
		buttonFloor++
		floorRequestButtonID++
	}
}

func (battery *Battery) createBasementFloorRequestButtons(amountOfBasements int) {
	buttonFloor := -1
	for i := 0; i < amountOfBasements; i++ {
		floorRequestButton := FloorRequestButton{
			ID:     floorRequestButtonID,
			status: "off",
			floor:  buttonFloor,
		}
		battery.floorRequestButtonsList = append(battery.floorRequestButtonsList, floorRequestButton)
		buttonFloor--
		floorRequestButtonID++
	}
}

func (battery *Battery) findBestColumn(requestedFloor int) Column {
	var bestColumn Column
	for _, column := range battery.columnsList {
		foundColumn := floorInColumn(requestedFloor, column.servedFloorsList)
		if foundColumn {
			bestColumn = column
		}
	}
	return bestColumn
}

func floorInColumn(x int, list []int) bool {
	for _, y := range list {
		if y == x {
			return true
		}
	}
	return false
}

func (battery *Battery) assignElevator(requestedFloor int, direction string) {
	fmt.Println("Passenger requests elevator at the lobby for floor", requestedFloor)
	var column Column = battery.findBestColumn(requestedFloor)
	fmt.Println(column.ID, "is the assigned column for this request")
	var elevator Elevator = column.findBestElevator(1, direction)
	fmt.Println(elevator.ID, "is the assigned elevator for this request")
	elevator.currentFloor = 1
	elevator.operateDoors()
	elevator.floorRequestList = append(elevator.floorRequestList, requestedFloor)
	//SORT FLOOR LIST
	fmt.Println("Elevator is moving")
	elevator.moveElevator()
	fmt.Println("Elevator is", elevator.status)
	elevator.operateDoors()
	if len(elevator.floorRequestList) == 0 {
		elevator.direction = ""
		elevator.status = "idle"
	}
	fmt.Println("Elevator is", elevator.status)

}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX//

//--------------------COLUMN CLASS--------------------//

type Column struct {
	ID                int
	status            string
	amountOfElevators int
	servedFloorsList  []int
	isBasement        bool
	elevatorsList     []Elevator
	callButtonsList   []CallButton
}

func (column *Column) createCallButtons(servedFloorsList []int, amountOfBasements int, isBasement bool) {
	callButtonID := 1
	if isBasement == true {
		buttonFloor := -1
		for i := 0; i < amountOfBasements; i++ {
			callButton := CallButton{
				ID:        callButtonID,
				status:    "off",
				floor:     buttonFloor,
				direction: "up",
			}
			column.callButtonsList = append(column.callButtonsList, callButton)
			buttonFloor--
			callButtonID++
		}
	} else {
		buttonFloor := 1
		for _, floor := range column.servedFloorsList {
			callButton := CallButton{
				ID:        callButtonID,
				status:    "off",
				floor:     floor,
				direction: "down",
			}
			column.callButtonsList = append(column.callButtonsList, callButton)
			buttonFloor++
			callButtonID++
		}

	}
}

func (column *Column) createElevators(servedFloorsList []int, amountOfElevators int) {
	elevatorID := 1
	for i := 0; i < amountOfElevators; i++ {
		elevator := Elevator{
			ID:           elevatorID,
			status:       "idle",
			servedFloors: column.servedFloorsList,
			currentFloor: 1,
		}
		column.elevatorsList = append(column.elevatorsList, elevator)
		elevatorID++
	}
}

func (column *Column) requestElevator(userFloor int, direction string) {
	fmt.Println("||Passenger requests elevator from", userFloor, "going", direction, "to the lobby||")
	var elevator Elevator = column.findBestElevator(userFloor, direction)
	fmt.Println("||", elevator.ID, "is the assigned elevator for this request||")
	elevator.floorRequestList = append(elevator.floorRequestList, userFloor)
	//elevator.sortFloorList()
	elevator.moveElevator()
	elevator.operateDoors()
}

func (column *Column) findBestElevator(floor int, direction string) Elevator {
	requestedFloor := floor
	requestedDirection := direction
	bestElevatorInfo := BestElevatorInfo{
		bestElevator: Elevator{},
		bestScore:    6,
		referenceGap: 1000000,
	}

	if requestedFloor == 1 {
		for _, elevator := range column.elevatorsList {
			if 1 == elevator.currentFloor && elevator.status == "stopped" {
				bestElevatorInfo = column.checkIfElevatorIsBetter(1, elevator, bestElevatorInfo, requestedFloor)
			} else if 1 == elevator.currentFloor && elevator.status == "idle" {
				bestElevatorInfo = column.checkIfElevatorIsBetter(2, elevator, bestElevatorInfo, requestedFloor)
			} else if 1 > elevator.currentFloor && elevator.direction == "up" {
				bestElevatorInfo = column.checkIfElevatorIsBetter(3, elevator, bestElevatorInfo, requestedFloor)
			} else if 1 < elevator.currentFloor && elevator.direction == "down" {
				bestElevatorInfo = column.checkIfElevatorIsBetter(3, elevator, bestElevatorInfo, requestedFloor)
			} else if elevator.status == "idle" {
				bestElevatorInfo = column.checkIfElevatorIsBetter(4, elevator, bestElevatorInfo, requestedFloor)
			} else {
				bestElevatorInfo = column.checkIfElevatorIsBetter(5, elevator, bestElevatorInfo, requestedFloor)
			}
		}
	} else {
		for _, elevator := range column.elevatorsList {
			if requestedFloor == elevator.currentFloor && elevator.status == "stopped" && requestedDirection == elevator.direction {
				bestElevatorInfo = column.checkIfElevatorIsBetter(1, elevator, bestElevatorInfo, requestedFloor)
			} else if requestedFloor > elevator.currentFloor && elevator.direction == "up" && requestedDirection == "up" {
				bestElevatorInfo = column.checkIfElevatorIsBetter(2, elevator, bestElevatorInfo, requestedFloor)
			} else if requestedFloor < elevator.currentFloor && elevator.direction == "down" && requestedDirection == "down" {
				bestElevatorInfo = column.checkIfElevatorIsBetter(2, elevator, bestElevatorInfo, requestedFloor)
			} else if elevator.status == "idle" {
				bestElevatorInfo = column.checkIfElevatorIsBetter(3, elevator, bestElevatorInfo, requestedFloor)
			} else {
				bestElevatorInfo = column.checkIfElevatorIsBetter(4, elevator, bestElevatorInfo, requestedFloor)
			}
		}
	}
	return bestElevatorInfo.bestElevator
}

func (column *Column) checkIfElevatorIsBetter(scoreToCheck int, newElevator Elevator, bestElevatorInfo BestElevatorInfo, floor int) BestElevatorInfo {

	if scoreToCheck < bestElevatorInfo.bestScore {
		bestElevatorInfo.bestScore = scoreToCheck
		bestElevatorInfo.bestElevator = newElevator
		bestElevatorInfo.referenceGap = int(math.Abs(float64(newElevator.currentFloor - floor)))
	} else if bestElevatorInfo.bestScore == scoreToCheck {
		gap := int(math.Abs(float64(newElevator.currentFloor - floor)))
		if bestElevatorInfo.referenceGap > gap {
			bestElevatorInfo.bestScore = scoreToCheck
			bestElevatorInfo.bestElevator = newElevator
			bestElevatorInfo.referenceGap = gap
		}
	}
	return bestElevatorInfo
}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX//

//--------------------ELEVATOR CLASS--------------------//

type Elevator struct {
	ID               int
	status           string
	servedFloors     []int
	currentFloor     int
	direction        string
	door             Door
	floorRequestList []int
}

func (elevator *Elevator) createDoors(doorID int, status string) {
	door := Door{
		ID:     doorID,
		status: status,
	}
	elevator.door = door
}

func (elevator *Elevator) moveElevator() {
	for i := 0; i < len(elevator.floorRequestList); i++ {
		destination := elevator.floorRequestList[0]
		elevator.status = "moving"
		if elevator.currentFloor < destination {
			elevator.direction = "Up"
			fmt.Println("Elevator ", elevator.ID, " status: ", elevator.status, " direction: ", elevator.direction)

			for i := elevator.currentFloor; i < destination; i++ {
				fmt.Println("Elevator ", elevator.ID, " status: ", elevator.status, " direction: ", elevator.direction, " current floor: ", elevator.currentFloor)
				elevator.currentFloor++
			}

		} else if elevator.currentFloor > destination {
			elevator.direction = "Down"
			fmt.Println("Elevator ", elevator.ID, " status: ", elevator.status, " direction: ", elevator.direction)

			for i := elevator.currentFloor; i > destination; i++ {
				fmt.Println("Elevator ", elevator.ID, " status: ", elevator.status, " direction: ", elevator.direction, " current floor: ", elevator.currentFloor)
				elevator.currentFloor--
			}
		}
		elevator.status = "stopped"
		fmt.Println("Elevator ", elevator.ID, " status: ", elevator.status, " current floor: ", elevator.currentFloor)
		elevator.floorRequestList = elevator.floorRequestList[1:]
	}

	elevator.status = "idle"
	fmt.Println("Elevator ", elevator.ID, " status: ", elevator.status)
	elevator.direction = "---"
	elevator.status = "idle"
	fmt.Println("Elevator ", elevator.ID, " status: ", elevator.status)
	elevator.direction = ""
}

func (elevator *Elevator) sortFloorList() {
	if elevator.direction == "Up" {
		sort.Ints(elevator.floorRequestList)
		fmt.Println(elevator.floorRequestList)

	} else if elevator.direction == "Down" {
		sort.Sort(sort.Reverse(sort.IntSlice(elevator.floorRequestList)))
		fmt.Println(elevator.floorRequestList)
	}
}

func (elevator *Elevator) operateDoors() {
	if elevator.status == "stopped" || elevator.status == "idle" {
		elevator.door.status = "open"
		fmt.Println("Doors: ", elevator.door.status)
		fmt.Println("(doors stay open for 6 seconds)")
		if len(elevator.floorRequestList) < 1 {
			elevator.direction = ""
			elevator.status = "idle"
			fmt.Println("Elevator ", elevator.ID, " status: ", elevator.status)
		}
	}
}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX//

//--------------------BEST ELEVATOR INFO--------------------//

type BestElevatorInfo struct {
	bestElevator Elevator
	bestScore    int
	referenceGap int
}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX//

//--------------------CALL BUTTON CLASS--------------------//

type CallButton struct {
	ID        int
	status    string
	floor     int
	direction string
}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX//

//--------------------FLOOR REQUEST BUTTON CLASS--------------------//

type FloorRequestButton struct {
	ID     int
	status string
	floor  int
}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX//

//--------------------DOOR CLASS--------------------//

type Door struct {
	ID     int
	status string
}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX//

//--------------------SCENARIO TESTS--------------------//

func scenario1() {
	battery := Battery{
		ID:                        1,
		amountOfColumns:           4,
		status:                    "online",
		amountOfFloors:            60,
		amountOfBasements:         6,
		amountOfElevatorPerColumn: 5,
		columnsList:               []Column{},
		floorRequestButtonsList:   []FloorRequestButton{},
	}
	if battery.amountOfBasements > 0 {
		battery.createBasementFloorRequestButtons(battery.amountOfBasements)
		battery.createBasementColumn(battery.amountOfBasements, battery.amountOfElevatorPerColumn)
		battery.amountOfColumns--
	}
	battery.createFloorRequestButtons(battery.amountOfFloors)
	battery.createColumns(battery.amountOfColumns, battery.amountOfFloors, battery.amountOfBasements, battery.amountOfElevatorPerColumn)
	battery.amountOfColumns++

	fmt.Println()
	fmt.Println("--------------------SCENARIO 1--------------------")
	fmt.Println()

	battery.columnsList[1].elevatorsList[0].currentFloor = 20
	battery.columnsList[1].elevatorsList[0].direction = "down"
	battery.columnsList[1].elevatorsList[0].status = "moving"
	battery.columnsList[1].elevatorsList[0].floorRequestList = append(battery.columnsList[1].elevatorsList[0].floorRequestList, 5)

	battery.columnsList[1].elevatorsList[1].currentFloor = 3
	battery.columnsList[1].elevatorsList[1].direction = "up"
	battery.columnsList[1].elevatorsList[1].status = "moving"
	battery.columnsList[1].elevatorsList[1].floorRequestList = append(battery.columnsList[1].elevatorsList[1].floorRequestList, 15)

	battery.columnsList[1].elevatorsList[2].currentFloor = 13
	battery.columnsList[1].elevatorsList[2].direction = "down"
	battery.columnsList[1].elevatorsList[2].status = "moving"
	battery.columnsList[1].elevatorsList[2].floorRequestList = append(battery.columnsList[1].elevatorsList[2].floorRequestList, 1)

	battery.columnsList[1].elevatorsList[3].currentFloor = 15
	battery.columnsList[1].elevatorsList[3].direction = "down"
	battery.columnsList[1].elevatorsList[3].status = "moving"
	battery.columnsList[1].elevatorsList[3].floorRequestList = append(battery.columnsList[1].elevatorsList[3].floorRequestList, 2)

	battery.columnsList[1].elevatorsList[4].currentFloor = 6
	battery.columnsList[1].elevatorsList[4].direction = "down"
	battery.columnsList[1].elevatorsList[4].status = "moving"
	battery.columnsList[1].elevatorsList[4].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

	battery.assignElevator(20, "up")
}

func scenario2() {
	battery := Battery{
		ID:                        1,
		amountOfColumns:           4,
		status:                    "online",
		amountOfFloors:            60,
		amountOfBasements:         6,
		amountOfElevatorPerColumn: 5,
		columnsList:               []Column{},
		floorRequestButtonsList:   []FloorRequestButton{},
	}
	if battery.amountOfBasements > 0 {
		battery.createBasementFloorRequestButtons(battery.amountOfBasements)
		battery.createBasementColumn(battery.amountOfBasements, battery.amountOfElevatorPerColumn)
		battery.amountOfColumns--
	}
	battery.createFloorRequestButtons(battery.amountOfFloors)
	battery.createColumns(battery.amountOfColumns, battery.amountOfFloors, battery.amountOfBasements, battery.amountOfElevatorPerColumn)
	battery.amountOfColumns++

	fmt.Println()
	fmt.Println("--------------------SCENARIO 2--------------------")
	fmt.Println()

	battery.columnsList[2].elevatorsList[0].currentFloor = 1
	battery.columnsList[2].elevatorsList[0].direction = "up"
	battery.columnsList[2].elevatorsList[0].status = "stopped"
	battery.columnsList[2].elevatorsList[0].floorRequestList = append(battery.columnsList[2].elevatorsList[0].floorRequestList, 21)

	battery.columnsList[2].elevatorsList[1].currentFloor = 23
	battery.columnsList[2].elevatorsList[1].direction = "up"
	battery.columnsList[2].elevatorsList[1].status = "moving"
	battery.columnsList[2].elevatorsList[1].floorRequestList = append(battery.columnsList[2].elevatorsList[1].floorRequestList, 28)

	battery.columnsList[2].elevatorsList[2].currentFloor = 33
	battery.columnsList[2].elevatorsList[2].direction = "down"
	battery.columnsList[2].elevatorsList[2].status = "moving"
	battery.columnsList[2].elevatorsList[2].floorRequestList = append(battery.columnsList[2].elevatorsList[2].floorRequestList, 1)

	battery.columnsList[2].elevatorsList[3].currentFloor = 40
	battery.columnsList[2].elevatorsList[3].direction = "down"
	battery.columnsList[2].elevatorsList[3].status = "moving"
	battery.columnsList[2].elevatorsList[3].floorRequestList = append(battery.columnsList[2].elevatorsList[3].floorRequestList, 24)

	battery.columnsList[2].elevatorsList[4].currentFloor = 39
	battery.columnsList[2].elevatorsList[4].direction = "down"
	battery.columnsList[2].elevatorsList[4].status = "moving"
	battery.columnsList[2].elevatorsList[4].floorRequestList = append(battery.columnsList[2].elevatorsList[4].floorRequestList, 1)

	battery.assignElevator(36, "up")
}

func scenario3() {
	battery := Battery{
		ID:                        1,
		amountOfColumns:           4,
		status:                    "online",
		amountOfFloors:            60,
		amountOfBasements:         6,
		amountOfElevatorPerColumn: 5,
		columnsList:               []Column{},
		floorRequestButtonsList:   []FloorRequestButton{},
	}
	if battery.amountOfBasements > 0 {
		battery.createBasementFloorRequestButtons(battery.amountOfBasements)
		battery.createBasementColumn(battery.amountOfBasements, battery.amountOfElevatorPerColumn)
		battery.amountOfColumns--
	}
	battery.createFloorRequestButtons(battery.amountOfFloors)
	battery.createColumns(battery.amountOfColumns, battery.amountOfFloors, battery.amountOfBasements, battery.amountOfElevatorPerColumn)
	battery.amountOfColumns++

	fmt.Println()
	fmt.Println("--------------------SCENARIO 3--------------------")
	fmt.Println()

	battery.columnsList[3].elevatorsList[0].currentFloor = 58
	battery.columnsList[3].elevatorsList[0].direction = "down"
	battery.columnsList[3].elevatorsList[0].status = "moving"
	battery.columnsList[3].elevatorsList[0].floorRequestList = append(battery.columnsList[3].elevatorsList[0].floorRequestList, 1)

	battery.columnsList[3].elevatorsList[1].currentFloor = 50
	battery.columnsList[3].elevatorsList[1].direction = "up"
	battery.columnsList[3].elevatorsList[1].status = "moving"
	battery.columnsList[3].elevatorsList[1].floorRequestList = append(battery.columnsList[3].elevatorsList[1].floorRequestList, 60)

	battery.columnsList[3].elevatorsList[2].currentFloor = 46
	battery.columnsList[3].elevatorsList[2].direction = "up"
	battery.columnsList[3].elevatorsList[2].status = "moving"
	battery.columnsList[3].elevatorsList[2].floorRequestList = append(battery.columnsList[3].elevatorsList[2].floorRequestList, 58)

	battery.columnsList[3].elevatorsList[3].currentFloor = 1
	battery.columnsList[3].elevatorsList[3].direction = "up"
	battery.columnsList[3].elevatorsList[3].status = "moving"
	battery.columnsList[3].elevatorsList[3].floorRequestList = append(battery.columnsList[3].elevatorsList[3].floorRequestList, 54)

	battery.columnsList[3].elevatorsList[4].currentFloor = 60
	battery.columnsList[3].elevatorsList[4].direction = "down"
	battery.columnsList[3].elevatorsList[4].status = "moving"
	battery.columnsList[3].elevatorsList[4].floorRequestList = append(battery.columnsList[3].elevatorsList[4].floorRequestList, 1)

	battery.columnsList[3].requestElevator(54, "down")
}

func scenario4() {
	battery := Battery{
		ID:                        1,
		amountOfColumns:           4,
		status:                    "online",
		amountOfFloors:            60,
		amountOfBasements:         6,
		amountOfElevatorPerColumn: 5,
		columnsList:               []Column{},
		floorRequestButtonsList:   []FloorRequestButton{},
	}
	if battery.amountOfBasements > 0 {
		battery.createBasementFloorRequestButtons(battery.amountOfBasements)
		battery.createBasementColumn(battery.amountOfBasements, battery.amountOfElevatorPerColumn)
		battery.amountOfColumns--
	}
	battery.createFloorRequestButtons(battery.amountOfFloors)
	battery.createColumns(battery.amountOfColumns, battery.amountOfFloors, battery.amountOfBasements, battery.amountOfElevatorPerColumn)
	battery.amountOfColumns++

	fmt.Println()
	fmt.Println("--------------------SCENARIO 4--------------------")
	fmt.Println()

	battery.columnsList[0].elevatorsList[0].currentFloor = -4

	battery.columnsList[0].elevatorsList[1].currentFloor = 1

	battery.columnsList[0].elevatorsList[2].currentFloor = -3
	battery.columnsList[0].elevatorsList[2].direction = "down"
	battery.columnsList[0].elevatorsList[2].status = "moving"
	battery.columnsList[0].elevatorsList[2].floorRequestList = append(battery.columnsList[0].elevatorsList[2].floorRequestList, -5)

	battery.columnsList[0].elevatorsList[3].currentFloor = -6
	battery.columnsList[0].elevatorsList[3].direction = "up"
	battery.columnsList[0].elevatorsList[3].status = "moving"
	battery.columnsList[0].elevatorsList[3].floorRequestList = append(battery.columnsList[0].elevatorsList[3].floorRequestList, 1)

	battery.columnsList[0].elevatorsList[4].currentFloor = -1
	battery.columnsList[0].elevatorsList[4].direction = "down"
	battery.columnsList[0].elevatorsList[4].status = "moving"
	battery.columnsList[0].elevatorsList[4].floorRequestList = append(battery.columnsList[0].elevatorsList[4].floorRequestList, -6)

	battery.columnsList[0].requestElevator(-3, "up")
}

//Instructions: TO SIMULATE A SCENARIO, SIMPLY UNCOMMENT THE DESIRED FUNCTION BY REMOVING THE "//" IN FRONT OF IT

func main() {
	//scenario1()
	//scenario2()
	//scenario3()
	//scenario4()
}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX//

package main

import (
	"fmt"
	"os"

	handlers "github.com/mannuR22/go-xlsx-data-processor.git/handlers"
	utility "github.com/mannuR22/go-xlsx-data-processor.git/utility"
)

func main() {
	//NOTE: In Go lang there is concept of SLICE instead of array which is more powerful alternative to array, however it functions similar to array
	// where ever i have stated array it actually means SLICE for increasing understandability of go developer

	var filePathToXLSX string = "./Assignment_Timecard.xlsx" /* path to xlsx file from root of project directory*/
	// Opening a file for writing. in case file doesn't exist ---> it creates new, otherwise ----> truncate existing file)
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error in creating file:", err.Error())
		return
	}
	//for successfully closing outputFile before termination of program
	defer outputFile.Close()

	//Redirecting standard output to the outputFile
	os.Stdout = outputFile

	//getting Map for (EmployeeName + FileNumber + PositionID) -> array of all records of same Employee
	// array of all records of same Employee will be sorted in ascending order according to Date (Time-In column)
	employeeToRecordsMap, err := utility.GetInfoMapFromXLSX(filePathToXLSX)
	if err != nil || employeeToRecordsMap == nil {
		fmt.Println("Error occurs while fetching info map from xlsx, Error:", err.Error())
		os.Exit(0)
	}

	//getting Array of Employees Info who worked for 'N' consecutive days (for 7 consecutive days N = 7) .
	ansA := handlers.NConsecutiveDaysEmployeesList(employeeToRecordsMap, 7 /* = N*/)
	utility.PrintInfo("Who has Worked for 7 consecutive days?", ansA)

	//getting Array of Employees Info who have less than maxHour = 10, but greater than minHour = 1 between shift
	ansB := handlers.HoursBetweenShiftEmployeesList(employeeToRecordsMap, 1 /* = minHour*/, 10 /* = maxHour*/)
	utility.PrintInfo("Who have less than 10 hours of time between shifts, but greater than 1 hour?", ansB)

	//getting Array of Employees Info who worked more than 14 hours, i.e.,minHours = 14
	ansC := handlers.WorkedMoreThanEmployeesList(employeeToRecordsMap, 14 /*minHours*/)
	//different Algorithm for filtering same form data make sure to comment out above line if you want to uncomment below line
	//ansC := handlers.WorkedMoreThanEmployeesList_V2(employeeToRecordsMap, 14 /*minHours*/)

	utility.PrintInfo("Who has Worked for more than 14hours in a single shift?", ansC)

	//changing standard output back to console
	os.Stdout = os.NewFile(0, "/dev/stdout")
}

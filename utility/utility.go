package utility

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/mannuR22/go-xlsx-data-processor.git/models"
)

func GetInfoMapFromXLSX(pathToXLSX string) (map[string][]models.Record, error) {
	//Opening the .xlsx file
	xlsx, err := excelize.OpenFile(pathToXLSX)
	if err != nil {
		fmt.Println("Error in getInfoMapFromXLSX utility method, while opening xlsx file.")
		return nil, err
	}

	//getting array of Sheets in .xlsx file
	sheetList := xlsx.GetSheetMap()

	//Checking sheetList should contain at least 1 sheet
	if len(sheetList) > 0 {

		//As we are working with the first and only sheet in .xlsx file
		sheetName := sheetList[1]

		//getting array of Row in sheet with name = sheetName
		rows := xlsx.GetRows(sheetName)

		//Allocating memory to map type string -> array Employee Records
		employeeToRecordsMap := make(map[string][]models.Record)

		// iterating over each row in rows
		for indx, row := range rows {

			//skipping iteration for header row (indx == 0) & record with missing data in any cell of row
			if indx == 0 || row[0] == "" || row[1] == "" || row[2] == "" ||
				row[3] == "" || row[4] == "" || row[5] == "" ||
				row[6] == "" || row[7] == "" || row[8] == "" {
				continue
			}

			timeTmp := strings.Split(string(row[4]), ":")
			hour, _ := strconv.Atoi(timeTmp[0])
			minute, _ := strconv.Atoi(timeTmp[1])

			//initializing custom type variable timeCard
			timeCard := models.TimeHoursIN{
				Hour:   hour,
				Minute: minute,
			}
			timeIn, err := excelNumberToTime(row[2])
			if err != nil {
				fmt.Println("Error in getInfoMapFromXLSX utility method, while converting TimeIn to time.")
				return nil, err
			}
			timeOut, err := excelNumberToTime(row[3])
			if err != nil {
				fmt.Println("Error in getInfoMapFromXLSX utility method, while converting TimeOut to time.")
				return nil, err
			}
			startPayCycle, err := excelNumberToDate(row[5])
			if err != nil {
				fmt.Println("Error in getInfoMapFromXLSX utility method, while converting startPayCycle to Date.")
				return nil, err
			}
			endPayCycle, err := excelNumberToDate(row[6])
			if err != nil {
				fmt.Println("Error in getInfoMapFromXLSX utility method, while converting endPayCycle to Date.")
				return nil, err
			}

			//Declaring and initializing custom type variable record
			record := models.Record{
				PositionID:        row[0],
				PositionStatus:    row[1],
				TimeIn:            timeIn,
				TimeOut:           timeOut,
				TimeHours:         timeCard,
				PayCycleStartDate: startPayCycle,
				PayCycleEndDate:   endPayCycle,
				EmployeeName:      row[7],
				FileNumber:        row[8],
			}

			//constructing unique key for making each employee unique
			mapKey := removeSpaces(record.PositionID) + removeSpaces(record.EmployeeName) + removeSpaces(record.FileNumber)
			if _, isExist := employeeToRecordsMap[mapKey]; isExist {
				//If Employee Already Exist in Map
				employeeToRecordsMap[mapKey] = append(employeeToRecordsMap[mapKey], record)
			} else {
				//If Employee Doesn't Exist in Map
				mapRecords := []models.Record{record}
				employeeToRecordsMap[mapKey] = mapRecords
			}

		}

		//sorting Records of Employee records by Time-In field/column
		for key, _ := range employeeToRecordsMap {
			sort.Slice(employeeToRecordsMap[key], func(i, j int) bool {
				return employeeToRecordsMap[key][i].TimeIn.Before(employeeToRecordsMap[key][j].TimeIn)
			})
		}

		return employeeToRecordsMap, nil

	}

	fmt.Println("utility,GetInfoMapFromXLSX: No Excel Sheet Found")
	return nil, nil
}

func GetDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func PrintInfo(q string, records []models.RecordOUT) {
	//Method for printing Employees Info
	fmt.Println("\n\nQ:", q)
	for indx, record := range records {

		fmt.Println("\nS.No:", indx+1)
		fmt.Println("Employee-Name:", record.Name)
		fmt.Println("Position-ID:", record.Position)
	}
}

func removeSpaces(input string) string {
	// Replace all white spaces with an empty string
	return strings.ReplaceAll(input, " ", "")
}

func excelNumberToDate(s string) (time.Time, error) {
	//parsing string to float type
	serialDate, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println("Error in utility:excelNumberToDate() while parsing to float.")
		return time.Time{}, err
	}

	// Excel uses a different reference date (January 1, 1900)
	referenceDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	// converting Serial date number to a time.Time value
	date := referenceDate.AddDate(0, 0, int(serialDate)-2) // Subtracting 2 to account for Excel's date offset
	return date, nil
}

func excelNumberToTime(s string) (time.Time, error) {
	//parsing string to float type
	serialNumber, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println("Error in utility:excelNumberToTime() while parsing to float.")
		return time.Time{}, err
	}

	// Excel uses a different reference date (January 1, 1900)
	referenceDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	// Calculating the days and fractional part
	days := int(serialNumber)
	fraction := serialNumber - float64(days)

	// Converting the serial number to a time.Time value
	date := referenceDate.AddDate(0, 0, days-2) // Subtracting 2 to account for Excel's date offset
	timeOfDay := time.Duration(fraction * 24 * float64(time.Hour))
	return date.Add(timeOfDay), nil
}

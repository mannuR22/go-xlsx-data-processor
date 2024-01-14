package handlers

import (
	"time"

	models "github.com/mannuR22/go-xlsx-data-processor.git/models"
	utility "github.com/mannuR22/go-xlsx-data-processor.git/utility"
)

func NConsecutiveDaysEmployeesList(nameToRecordsMap map[string][]models.Record, N int) []models.RecordOUT {

	//declaring empty array of reqRecords of custom type RecordOUT for saving required employee info
	reqRecords := []models.RecordOUT{}

	//iterating over each unique employee records
	for _, records := range nameToRecordsMap {

		if len(records) == 0 {
			continue
		}

		//Checking Max N Consecutive Days As records is already sorted so consecutive Time-In Dates will be present adjacent to each other
		var lastDate time.Time
		var consecutiveDays int = 0
		var isFound bool = false

		//iterating over each records belongs to unique employee
		for indx, record := range records {

			//extracting date from TimeIn column data of record
			currDate := utility.GetDate(record.TimeIn)

			if indx == 0 {
				lastDate = currDate
				consecutiveDays = 1
				continue
			}

			//calculating days difference between last record and current record's TimeIN date
			difference := int(currDate.Sub(lastDate).Hours() / 24)

			if difference == 0 {
				//if record belongs to shift of same day then skip current iteration
				continue
			} else if difference == 1 {
				//if record is of  next consecutive day
				consecutiveDays++
			} else {

				if consecutiveDays >= N {
					reqRecords = append(reqRecords, models.RecordOUT{
						Name:     record.EmployeeName,
						Position: record.PositionID,
					})
					isFound = true
					break
				}
				//resetting value of consecutive days count if last consecutive days count is not equals to N
				consecutiveDays = 1

			}

			//updating lastDate to current Date of current iteration
			lastDate = currDate

		}

		if !isFound && consecutiveDays >= N {
			//In case employee is consecutive for all days in records & more than or equal to consecutiveDay N
			reqRecords = append(reqRecords, models.RecordOUT{
				Name:     records[0].EmployeeName,
				Position: records[0].PositionID,
			})
		}

	}

	return reqRecords
}

func HoursBetweenShiftEmployeesList(nameToRecordsMap map[string][]models.Record, minHour, maxHour int) []models.RecordOUT {

	//declaring empty array of reqRecords of custom type RecordOUT for saving required employee info
	reqRecords := []models.RecordOUT{}

	//iterating over each unique employee records
	for _, records := range nameToRecordsMap {

		//array of custom type Record to store shifts belongs to same day
		shifts := []models.Record{}
		var isFound bool = false

		//iterating over each records belongs to unique employee
		for _, record := range records {

			//appending record which belong to same day
			if len(shifts) == 0 || utility.GetDate(shifts[0].TimeIn).Unix() == utility.GetDate(record.TimeIn).Unix() {
				shifts = append(shifts, record)
				continue
			} else {

				//analyzing shift array if there is change in date of current shift.

				for i := 1; i < len(shifts); i++ {
					difference := int(shifts[i].TimeIn.Sub(shifts[i-1].TimeOut).Hours())
					if difference > minHour && difference < maxHour {
						reqRecords = append(reqRecords, models.RecordOUT{
							Name:     record.EmployeeName,
							Position: record.PositionID,
						})
						isFound = true
						break
					}
				}

				//updating shift array with current record as previous day shifts data is already processed
				shifts = []models.Record{record}

				if isFound {
					break
				}
			}
		}

		//For handling edge case i.e., when few shifts from End of Records fall under same day
		for i := 1; !isFound && i < len(shifts); i++ {
			difference := int(shifts[i].TimeIn.Sub(shifts[i-1].TimeOut).Hours())
			if difference > minHour && difference < maxHour {
				reqRecords = append(reqRecords, models.RecordOUT{
					Name:     records[0].EmployeeName,
					Position: records[0].PositionID,
				})
				break
			}
		}
	}

	return reqRecords
}

func WorkedMoreThanEmployeesList(nameToRecordsMap map[string][]models.Record, minHours int) []models.RecordOUT {

	//declaring empty array of reqRecords of custom type RecordOUT for saving required employee info
	reqRecords := []models.RecordOUT{}

	//iterating over each unique employee records
	for _, records := range nameToRecordsMap {

		// iterating over each record of unique employee
		for _, record := range records {

			//checking shift Hours using using Time Hour Card column in xlsx
			if record.TimeHours.Hour > minHours ||
				(record.TimeHours.Hour == minHours && record.TimeHours.Minute > 0) {
				reqRecords = append(reqRecords, models.RecordOUT{
					Name:     record.EmployeeName,
					Position: record.PositionID,
				})
				break
			}
		}
	}

	return reqRecords

}

func WorkedMoreThanEmployeesList_V2(nameToRecordsMap map[string][]models.Record, minHours int) []models.RecordOUT {

	//declaring empty array of reqRecords of custom type RecordOUT for saving required employee info
	reqRecords := []models.RecordOUT{}

	//iterating over each unique employee records
	for _, records := range nameToRecordsMap {

		// iterating over each record of unique employee
		for _, record := range records {

			//calculating hours difference between TimeIn and TimeOut column
			difference := int(record.TimeOut.Sub(record.TimeIn).Hours())

			if difference > minHours {
				//if difference is more than minHours then it will add the record to reqRecord
				// and terminate loop for current employee
				reqRecords = append(reqRecords, models.RecordOUT{
					Name:     record.EmployeeName,
					Position: record.PositionID,
				})
				break
			}
		}
	}

	return reqRecords

}

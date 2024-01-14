package models

import "time"

// Custom type structure for storing Time Hour Card column / TimeHours Filed for Record Struct
type TimeHoursIN struct {
	Hour   int
	Minute int
}

// Custom Type structure for storing each record for employee
type Record struct {
	PositionID        string
	PositionStatus    string
	TimeIn            time.Time
	TimeOut           time.Time
	TimeHours         TimeHoursIN
	PayCycleStartDate time.Time
	PayCycleEndDate   time.Time
	EmployeeName      string
	FileNumber        string
}

// Custom Type structure for storing require Employee Info
type RecordOUT struct {
	Position string
	Name     string
}

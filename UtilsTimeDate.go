package Utils

import (
	"time"
)

//////////////////////////////////////////////////////

var UDateTime _DateTime_s
type _DateTime_s struct {
	/*
		GetDateTimeStr gets the current time and date in the format DATE_TIME_FORMAT.

		-----------------------------------------------------------

		> Returns:
		  - the current time and date in the default format
	*/
	GetDateTimeStr func() string
	/*
		GetTimeStr gets the current time in the format TIME_FORMAT.

		-----------------------------------------------------------

		> Returns:
		  - the current time in the default format
	*/
	GetTimeStr func() string
	/*
		GetDateStr gets the current date in the format DATE_FORMAT.

		-----------------------------------------------------------

		> Returns:
		  - the current date in the default format
	*/
	GetDateStr func() string
}
//////////////////////////////////////////////////////

const TIME_FORMAT string = "15:04:05"
const DATE_FORMAT string = "2006-01-02"
const DATE_TIME_FORMAT string = DATE_FORMAT + " -- " + TIME_FORMAT + " (MST)"

func getDateTimeStrTIMEDATE() string {
	return time.Now().Format(DATE_TIME_FORMAT)
}

func getDateStrTIMEDATE() string {
	return time.Now().Format(DATE_FORMAT)
}

func getTimeStrTIMEDATE() string {
	return time.Now().Format(TIME_FORMAT)
}

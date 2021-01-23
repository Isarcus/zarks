package system

import (
	"strconv"
	"time"
)

// GetTimeAsString returns the current system time formatted as "[year.month.day] [hour.minute.second]"
func GetTimeAsString(t time.Time) string {
	str := strconv.Itoa
	bigTime := "[" + str(t.Year()) + "." + str(int(t.Month())) + "." + str(t.Day()) + "]"
	smallTime := "[" + str(t.Hour()) + "." + str(t.Minute()) + "." + str(t.Second()) + "]"
	return bigTime + " " + smallTime
}

package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func LookupEnvOrString(key string, defaultValue string) string {
	envVariable, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return envVariable
}

func LookupEnvOrInt(key string, defaultValue int) int {
	envVariable, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	num, err := strconv.Atoi(envVariable)
	if err != nil {
		panic(err.Error())
	}
	return num
}

func IsValidTime(time string) bool {
	/*
	   Use regex to match military time in <HH>:<MM>
	   source: https://stackoverflow.com/questions/1494671/regular-expression-for-matching-time-in-military-24-hour-format
	   regex: ^([01]\d|2[0-3]):([0-5]\d)$
	*/
	matched, _ := regexp.MatchString(`^([01]\d|2[0-3]):([0-5]\d)$`, time)
	return matched
}

func ParseDayOfMonth(day int) string {
	/*
		day is 1-31
		1 -> 1st
		2 -> 2nd
		3 -> 3rd
		4 -> 4th
		...
	*/
	ones_digit := day % 10
	if ones_digit == 1 {
		return fmt.Sprintf("%vst", day)
	} else if ones_digit == 2 {
		return fmt.Sprintf("%vnd", day)
	} else if ones_digit == 3 {
		return fmt.Sprintf("%vrd", day)
	} else {
		return fmt.Sprintf("%vth", day)
	}
}

// https://stackoverflow.com/questions/73880828/list-the-number-of-days-in-current-date-month
func DaysInMonth(t time.Time) int {
	t = time.Date(t.Year(), t.Month(), 32, 0, 0, 0, 0, time.UTC)
	daysInMonth := 32 - t.Day()
	return daysInMonth
}

func ParseReminderTime(reminderTime string) (int, int) {
	t := strings.Split(reminderTime, ":")
	hour, _ := strconv.Atoi(t[0])
	minute, _ := strconv.Atoi(t[1])
	return hour, minute
}

package types

import (
	"log"
	"strings"
	"time"
)

func DBIntToBool(i int) bool {
	if i == 1 {
		return true
	}
	return false
}

type DAYS_PER_WEEK string

const (
	SUNDAY    DAYS_PER_WEEK = "sunday"
	MONDAY    DAYS_PER_WEEK = "monday"
	TUESDAY   DAYS_PER_WEEK = "tuesday"
	WEDNESDAY DAYS_PER_WEEK = "wednesday"
	THURSDAY  DAYS_PER_WEEK = "thursday"
	FRIDAY    DAYS_PER_WEEK = "friday"
	SATURDAY  DAYS_PER_WEEK = "saturday"
)

func GetTodaysDayName() DAYS_PER_WEEK {
	return ValidateDayOfTheWeek(time.Now().Weekday().String())
}

// i can assert the type... but i prefer to make sure...
func ValidateDayOfTheWeek(d string) DAYS_PER_WEEK {
	switch strings.TrimSpace(strings.ToLower(d)) {
	case string(SUNDAY):
		return SUNDAY
	case string(MONDAY):
		return MONDAY
	case string(TUESDAY):
		return TUESDAY
	case string(WEDNESDAY):
		return WEDNESDAY
	case string(THURSDAY):
		return THURSDAY
	case string(FRIDAY):
		return FRIDAY
	case string(SATURDAY):
		return SATURDAY
	default:
		log.Fatalf("Invalid day name: %s. Remove or update it from the config file.\n", d)
	}

	return ""
}

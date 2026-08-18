package types

import (
	"log"
	"strings"
	"time"
)

type ALL_MONTHS string

const (
	JANUARY   ALL_MONTHS = "january"
	FEBRUARY  ALL_MONTHS = "february"
	MARCH     ALL_MONTHS = "march"
	APRIL     ALL_MONTHS = "april"
	MAY       ALL_MONTHS = "may"
	JUNE      ALL_MONTHS = "june"
	JULY      ALL_MONTHS = "july"
	AUGUST    ALL_MONTHS = "august"
	SEPTEMBER ALL_MONTHS = "september"
	OCTOBER   ALL_MONTHS = "october"
	NOVEMBER  ALL_MONTHS = "november"
	DECEMBER  ALL_MONTHS = "december"

	EVERY ALL_MONTHS = "every"
)

type DATE_FORMATS string

const (
	MM_DD_YYYY DATE_FORMATS = "01/02/2006" // <- only support this for now
	//DD_MM_YYYY = "02/01/2006"
)

func GetDate(format DATE_FORMATS) string { return time.Now().Format(string(format)) }
func GetTodaysMonth() ALL_MONTHS         { return ValidateMonth(time.Now().Month().String()) }

// i can assert the type... but i prefer to make sure...
func ValidateMonth(d string) ALL_MONTHS {
	m := strings.TrimSpace(strings.ToLower(d))

	switch m {
	case string(JANUARY):
	case string(FEBRUARY):
	case string(MARCH):
	case string(APRIL):
	case string(MAY):
	case string(JUNE):
	case string(JULY):
	case string(AUGUST):
	case string(SEPTEMBER):
	case string(OCTOBER):
	case string(NOVEMBER):
	case string(DECEMBER):
	case string(EVERY):
	default:
		log.Fatalf("Invalid day name: %s. Remove or update it from the config file.\n", d)
	}

	return ALL_MONTHS(m)
}

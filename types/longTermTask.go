package types

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"unicode"
)

type PRIORITY_TYPES string

const (
	PRIO_LOW  PRIORITY_TYPES = "low"
	PRIO_MED  PRIORITY_TYPES = "med"
	PRIO_HIGH PRIORITY_TYPES = "high"
)

type LongTermTasksCfg struct {
	Ends     string         `toml:"ends"`
	Name     string         `toml:"name"`
	Priority PRIORITY_TYPES `toml:"priority"`

	//custom
	MM_DD_YYYY_DATE string
}

type EXPIRES_NOTATION string

const (
	EXPIRES_IN_DAYS_NOTATION  EXPIRES_NOTATION = "d"
	EXPIRES_IN_WEEK_NOTATION  EXPIRES_NOTATION = "w"
	EXPIRES_IN_MONTH_NOTATION EXPIRES_NOTATION = "m"

	MAX_DATE_NUMBER_ALLOWED int = 999
)

var (
	expiration_example_of_usage = fmt.Sprintf(`
	Usage (two ways): 
		- mm/dd/yy (month, day, year) [RECOMMENDED]
		- {date_number 1-%d}{date_notation (%s, %s ,%s)}
	For example:
		- 120%s
		- 57%s
		- 1%s
		- 20/04/2030
		- 31/12/2026
	`, MAX_DATE_NUMBER_ALLOWED,
		EXPIRES_IN_DAYS_NOTATION, EXPIRES_IN_WEEK_NOTATION, EXPIRES_IN_MONTH_NOTATION,
		EXPIRES_IN_DAYS_NOTATION, EXPIRES_IN_WEEK_NOTATION, EXPIRES_IN_MONTH_NOTATION)
)

// user can insert two different ways of date:
// 31/12/2026 (datestamp)
// 7d (day notation)
func (l *LongTermTasksCfg) ValidateExpires_In() error {
	e := l.Ends
	t, err := time.Parse(e, string(MM_DD_YYYY))
	now := time.Now()

	//NOTE: if the datestamp format is valid:
	if err == nil {
		if t.Compare(now) != -1 {
			return fmt.Errorf("We're not in the past. Select a valid date that's %s or further. \n%s", now.Format(string(MM_DD_YYYY)), expiration_example_of_usage)
		}

		if (now.Sub(t).Hours() / 24) > float64(MAX_DATE_NUMBER_ALLOWED) {
			return fmt.Errorf("date_number is higher than %d. \n%s", MAX_DATE_NUMBER_ALLOWED, expiration_example_of_usage)
		}

		l.MM_DD_YYYY_DATE = e
		return nil
	}

	//NOTE: else, it should be day notation
	var notation_idx int = -1
	for i, r := range e {
		if unicode.IsLetter(r) {
			if notation_idx != -1 {
				return fmt.Errorf("Contains more than one date_notation. \n%s", expiration_example_of_usage)
			}
			notation_idx = i
		}
	}

	if notation_idx == -1 || notation_idx == 0 {
		return fmt.Errorf("Needs to contain at least one date_notation and one date_number. \n%s", expiration_example_of_usage)
	}

	notation := EXPIRES_NOTATION(e[notation_idx])
	if err := l.validateNotation(notation); err != nil {
		return fmt.Errorf("%s. \n%s", err.Error(), expiration_example_of_usage)
	}

	number_of_exp, err := strconv.Atoi(e[:notation_idx])
	if err != nil {
		return fmt.Errorf("Couldn't convert date_number string to int. \n%s", expiration_example_of_usage)
	}

	if number_of_exp > MAX_DATE_NUMBER_ALLOWED {
		return fmt.Errorf("date_number is higher than %d. \n%s", MAX_DATE_NUMBER_ALLOWED, expiration_example_of_usage)
	}

	a := now.AddDate(0, 0,
		l.sumOnDifferentDayNotations(notation, number_of_exp)/24)
	l.MM_DD_YYYY_DATE = a.Format(string(MM_DD_YYYY))

	return nil
}

func (l *LongTermTasksCfg) ValidatePriority() error {
	switch l.Priority {
	case PRIO_LOW:
		return nil
	case PRIO_MED:
		return nil
	case PRIO_HIGH:
		return nil
	default:
		return fmt.Errorf(
			`Priority '%s' is not a valid priority of a Long Term Task.
			Did you meant one of the following:
			- "%s"
			- "%s"
			- "%s"

			`, string(l.Priority), PRIO_LOW, PRIO_MED, PRIO_HIGH)
	}
}

// WARN: PRIVATE
func (_ *LongTermTasksCfg) validateNotation(n EXPIRES_NOTATION) error {
	switch n {
	case EXPIRES_IN_DAYS_NOTATION:
		return nil
	case EXPIRES_IN_WEEK_NOTATION:
		return nil
	case EXPIRES_IN_MONTH_NOTATION:
		return nil
	default:
		return fmt.Errorf("Notation %s is not valid", n)
	}
}

const (
	HOURS_IN_A_DAY  = 24
	HOURS_IN_A_WEEK = 168
)

// one day = 24hs
func (l *LongTermTasksCfg) sumOnDifferentDayNotations(n EXPIRES_NOTATION, count int) (total_hours int) {
	total_hours = 0

	switch n {
	case EXPIRES_IN_DAYS_NOTATION:
		total_hours = count * HOURS_IN_A_DAY
	case EXPIRES_IN_WEEK_NOTATION:
		total_hours = count * HOURS_IN_A_WEEK
	case EXPIRES_IN_MONTH_NOTATION:
		total_hours = l.daysInMonths(count) * HOURS_IN_A_DAY
	default:
		log.Fatalln("Invalid notation: ", n)
	}

	return
}

func (_ *LongTermTasksCfg) daysInMonths(how_many_months int) int {
	n := time.Now()
	var days_in_total int = 0

	for i := range how_many_months {
		m := n.AddDate(0, i+1, 0)
		days_in_total += time.Date(
			m.Year(),
			m.Month()+1,
			0, 0, 0, 0, 0, time.UTC).Day()
	}
	return days_in_total
}

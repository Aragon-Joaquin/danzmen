package types

import (
	"fmt"
	"strconv"
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
	Usage: 
		- {date_number 1-%d}{date_notation (%s, %s ,%s)}
	For example:
		- 120%s
		- 57%s
		- 1%s
	`, MAX_DATE_NUMBER_ALLOWED,
		EXPIRES_IN_DAYS_NOTATION, EXPIRES_IN_WEEK_NOTATION, EXPIRES_IN_MONTH_NOTATION,
		EXPIRES_IN_DAYS_NOTATION, EXPIRES_IN_WEEK_NOTATION, EXPIRES_IN_MONTH_NOTATION)
)

func (l *LongTermTasksCfg) ValidateExpires_In() (EXPIRES_NOTATION, uint32, error) {
	var def_notation EXPIRES_NOTATION = EXPIRES_IN_DAYS_NOTATION
	var def_date_numb uint32 = 0

	e := l.Ends
	var notation_idx int = -1
	for i, r := range e {
		if unicode.IsLetter(r) {
			if notation_idx != -1 {
				return def_notation, def_date_numb, fmt.Errorf("Contains more than one date_notation. \n%s", expiration_example_of_usage)
			}
			notation_idx = i
		}
	}

	if notation_idx == -1 || notation_idx == 0 {
		return def_notation, def_date_numb, fmt.Errorf("Needs to contain at least one date_notation and one date_number. \n%s", expiration_example_of_usage)
	}

	notation := EXPIRES_NOTATION(e[notation_idx])
	if err := l.validateNotation(notation); err != nil {
		return def_notation, def_date_numb, fmt.Errorf("%s. \n%s", err.Error(), expiration_example_of_usage)
	}

	number_of_exp, err := strconv.Atoi(e[:notation_idx])
	if err != nil {
		return def_notation, def_date_numb, fmt.Errorf("Couldn't convert date_number string to int. \n%s", expiration_example_of_usage)
	}

	if number_of_exp > MAX_DATE_NUMBER_ALLOWED {
		return def_notation, def_date_numb, fmt.Errorf("date_number is higher than %d. \n%s", MAX_DATE_NUMBER_ALLOWED, expiration_example_of_usage)
	}

	return notation, uint32(number_of_exp), nil
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

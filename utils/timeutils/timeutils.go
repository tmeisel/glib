package timeutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errPkg "github.com/tmeisel/glib/error"
	"github.com/tmeisel/glib/utils/strutils"
)

// ParseDate parses a string like 20060102 and returns it as int(2006), time.January, int(2).
// In case of an error, the returned date is set to January 1st, 0001 (aka the zero time's date)
func ParseDate(date string) (year int, month time.Month, day int, err error) {
	if len(date) != 8 {
		return 1, time.January, 1, errPkg.NewUserMsg(nil, "invalid format. expecting YYYYMMDD")
	}

	yearStr, monthStr, dayStr :=
		strutils.SubString(date, 0, 4),
		strutils.SubString(date, 4, 2),
		strutils.SubString(date, 6, 2)

	year, err = strconv.Atoi(yearStr)
	if err != nil {
		return 1, time.January, 1, err
	}

	var monthInt int
	monthInt, err = strconv.Atoi(monthStr)
	if err != nil {
		return 1, time.January, 1, err
	}

	month = time.Month(monthInt)
	if time.January > month || time.December < month {
		return 1, time.January, 1, errPkg.NewUserMsg(nil, "invalid month")
	}

	day, err = strconv.Atoi(dayStr)
	if err != nil {
		return 1, time.January, 1, err
	}

	return year, time.Month(monthInt), day, nil
}

// AddDate takes a duration as a string in EITHER
// the form "3y" (years), "3M" (months) or "3d" (days)
// and adds it to the given time.Time to. Any other
// specification or combination will throw an error
func AddDate(to time.Time, dur string) (time.Time, error) {
	dur = strings.TrimSpace(dur)

	if len(dur) == 0 {
		return to, errPkg.NewUserMsg(nil, "invalid duration string")
	}

	switch strutils.SubString(dur, 0, -1) {
	case "y":
		yearsStr := strutils.SubString(dur, 0, len(dur)-1)
		years, err := strconv.Atoi(yearsStr)
		if err != nil {
			return to, errPkg.NewUserMsg(err, fmt.Sprintf("invalid duration string '%s' specified", dur))
		}

		return to.AddDate(years, 0, 0), nil
	case "M":
		monthsStr := strutils.SubString(dur, 0, len(dur)-1)
		months, err := strconv.Atoi(monthsStr)
		if err != nil {
			return to, errPkg.NewUserMsg(err, fmt.Sprintf("invalid duration string '%s' specified", dur))
		}

		return to.AddDate(0, months, 0), nil
	case "d":
		daysStr := strutils.SubString(dur, 0, len(dur)-1)
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return to, errPkg.NewUserMsg(err, fmt.Sprintf("invalid duration string '%s' specified", dur))
		}

		return to.AddDate(0, 0, days), nil
	}

	return to, errPkg.NewUserMsg(nil, fmt.Sprintf("invalid duration string '%s' specified", dur))
}

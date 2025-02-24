package timeutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	errPkg "github.com/tmeisel/glib/error"
	"github.com/tmeisel/glib/utils/strutils"
)

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

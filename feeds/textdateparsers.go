package feeds

import (
	"strconv"
	"strings"
	"time"
)

type TextDateParser func(textDate string) (date time.Time, err error)

func DashedYearDayMonthDateParser() TextDateParser {
	return func(textDate string) (date time.Time, err error) {
		return dashedYearDayMonthDateParser(textDate, time.UTC)
	}
}

func DashedYearDayMonthDateParserForLocation(location *time.Location) TextDateParser {
	return func(textDate string) (date time.Time, err error) {
		return dashedYearDayMonthDateParser(textDate, location)
	}
}

func dashedYearDayMonthDateParser(textDate string, location *time.Location) (date time.Time, err error) {
	splits := strings.Split(textDate, "-")
	year, err := strconv.ParseInt(splits[0], 10, 0)
	month, err := strconv.ParseInt(splits[1], 10, 0)
	day, err := strconv.ParseInt(splits[2], 10, 0)
	date = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, location)

	return date, err
}

func slashedYearDayMonthDateParser(textDate string) (date time.Time, err error) {
	splits := strings.Split(textDate, "/")
	year, err := strconv.ParseInt(splits[0], 10, 0)
	month, err := strconv.ParseInt(splits[1], 10, 0)
	day, err := strconv.ParseInt(splits[2], 10, 0)
	date = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)

	return date, err
}

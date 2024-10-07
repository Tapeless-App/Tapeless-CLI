package util

import (
	"errors"
	"time"
)

func DateTimeToDateStr(input string) string {

	if len(input) == 0 {
		return ""
	}

	date, err := time.Parse("2006-01-02T15:04:05.000Z", input)

	if err != nil {
		return ""
	}

	return date.Format("2006-01-02")

}

func VerifyStartBeforeEnd(start string, end string, layout string) error {
	startDate, err := time.Parse(layout, start)

	if err != nil {
		return err
	}

	if end == "" {
		return nil
	}

	endDate, err := time.Parse(layout, end)

	if err != nil {
		return err
	}

	if startDate.After(endDate) {
		return errors.New("start date must be before end date")
	}

	return nil
}

func IsDateInPast(layout string, date string) (bool, error) {

	if date == "" {
		return false, errors.New("no date provided")
	}

	parsedDate, err := time.Parse(layout, date)

	if err != nil {
		println("Error parsing:", date, err)
		return false, err
	}

	if time.Now().After(parsedDate) {
		return true, nil
	}

	return false, nil

}

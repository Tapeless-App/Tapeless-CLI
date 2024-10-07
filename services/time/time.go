package timeService

import (
	"fmt"

	"tapeless.app/tapeless-cli/env"
	"tapeless.app/tapeless-cli/util"
)

type TimeEntryCreateRequest struct {
	Date        string  `json:"date"`
	Hours       float64 `json:"hours"`
	Description string  `json:"description"`
}

type TimeEntryCreateResponse struct {
	TimeEntriesCount int     `json:"timeEntriesCount"`
	TotalHours       float64 `json:"totalHours"`
}

type TimeEntry struct {
	Id          int     `json:"id"`
	Date        string  `json:"date"`
	Hours       float64 `json:"hours"`
	Description string  `json:"description"`
}

type TimeEntryFetchResponse = []TimeEntry

func CreateTimeEntry(projectId int, request TimeEntryCreateRequest) (TimeEntryCreateResponse, error) {

	timeCreatereponse := TimeEntryCreateResponse{}

	err := util.MakeAuthRequestAndParseResponse("POST", fmt.Sprintf("%s/projects/%d/time", env.ApiURL, projectId),
		request, &timeCreatereponse)

	return timeCreatereponse, err

}

func FetchTimeEntries(projectId int, date string) (TimeEntryFetchResponse, error) {

	timeEntriesResponse := TimeEntryFetchResponse{}

	err := util.MakeAuthRequestAndParseResponse("GET", fmt.Sprintf("%s/projects/%d/time?date=%s", env.ApiURL, projectId, date),
		nil, &timeEntriesResponse)

	return timeEntriesResponse, err

}

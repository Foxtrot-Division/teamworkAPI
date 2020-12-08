package teamworkapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// TimeEntry models an individual time entry.
type TimeEntry struct {
	ID       string    `json:"id"`
	PersonID string    `json:"person-id"`
	Hours    string    `json:"hours"`
	Minutes  string    `json:"minutes"`
	Date     time.Time `json:"date"`
}

// TimeEntries models an array of time entries.
type TimeEntries struct {
	TimeEntries []TimeEntry `json:"time-entries"`
}

func (conn connection) GetTimeEntriesByPerson(personID string, from string, to string) (*TimeEntries, error) {

	queryParams := make(map[string]interface{})

	if personID != "" {
		queryParams["userId"] = personID
	}

	if from != "" {
		_, err := time.Parse("20060102", from)
		if err != nil {
			return nil, fmt.Errorf("invalid format for from parameter.  Should be YYYYMMDD, but found %s", from)
		}
		queryParams["fromdate"] = from
	}

	if to != "" {
		_, err := time.Parse("20060102", to)
		if err != nil {
			return nil, fmt.Errorf("invalid format for to parameter.  Should be YYYYMMDD, but found %s", from)
		}
		queryParams["todate"] = to
	}

	data, err := conn.GetRequest("time_entries", queryParams)

	if err != nil {
		return nil, err
	}

	t := new(TimeEntries)

	err = json.Unmarshal(data, &t)

	if err != nil {
		return nil, err
	}

	return t, nil
}

// SumHours returns the total hours for a specified user found in the TimeEntries array.
func (e *TimeEntries) SumHours(personID string) (float64, error) {
	found := false
	hours := 0.0

	for _, v := range e.TimeEntries {
		if v.PersonID == personID {
			h, err := strconv.ParseFloat(v.Hours, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to convert %s to float64", v.Hours)
			}

			found = true
			hours += h

			if v.Minutes != "0" {
				m, err := strconv.ParseFloat(v.Minutes, 64)
				if err != nil {
					return 0, fmt.Errorf("failed to convert %s to float64", v.Minutes)
				}

				hours += m / 60
			}
		}
	}

	if !found {
		return 0, fmt.Errorf("user ID %s not found in response data", personID)
	}

	return hours, nil
}

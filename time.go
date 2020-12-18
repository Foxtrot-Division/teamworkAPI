package teamworkapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// TimeEntry models an individual time entry.
type TimeEntry struct {
	ID          string   `json:"id"`
	PersonID    string   `json:"person-id"`
	Description string   `json:"description"`
	Hours       string   `json:"hours"`
	Minutes     string   `json:"minutes"`
	Date        string   `json:"date"`			// expected format is YYYYMMDD
	IsBillable  string   `json:"isbillable"`
	ProjectID   string   `json:"project-id"`
	TaskID      string   `json:"todo-item-id"`
}

// TimeEntryJSON provides a wrapper around TimeEntry to properly marshal json
// data when posting to API.
type TimeEntryJSON struct {
	Entry *TimeEntry `json:"time-entry"`
}

// TimeEntries models an array of time entries.
type TimeEntries struct {
	TimeEntries []TimeEntry `json:"time-entries"`
}

// GetTimeEntriesByPerson retrieves time entries for a specific Teamwork user, for the specified time period.
func (conn Connection) GetTimeEntriesByPerson(personID string, from string, to string) (*TimeEntries, error) {
	
	errBuff := ""

	if personID == "" {
		errBuff += "personID"
	} else {
		_, err := strconv.Atoi(personID)
		if err != nil {
			return nil, err
		}
	}

	if from == "" {
		if errBuff != "" {
			errBuff += ", "
		}
		errBuff += "from"
	}

	if to == "" {
		if errBuff != "" {
			errBuff += ", "
		}
		errBuff += "to"
	}

	if errBuff != "" {
		return nil, fmt.Errorf("missing required parameter(s): %s", errBuff)
	}

	queryParams := make(map[string]interface{})


	queryParams["userId"] = personID
	
	_, err := time.Parse("20060102", from)
	if err != nil {
		return nil, fmt.Errorf("invalid format for from parameter.  Should be YYYYMMDD, but found %s", from)
	}
	queryParams["fromdate"] = from

	_, err = time.Parse("20060102", to)
	if err != nil {
		return nil, fmt.Errorf("invalid format for to parameter.  Should be YYYYMMDD, but found %s", from)
	}
	queryParams["todate"] = to

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

// PostTimeEntry posts an individual time entry to the specified task.  The time
// entry is posted to the task ID found in the entry parameter.
func (conn *Connection) PostTimeEntry(entry *TimeEntry) (string, error) {

	errBuff := ""

	if entry.PersonID == "" {
		errBuff += "PersonID"
	}

	if entry.TaskID == "" {
		if errBuff != "" {
			errBuff += ", "
		}
		errBuff += "TaskID"
	}

	if entry.Date == "" {
		if errBuff != "" {
			errBuff += ", "
		}
		errBuff += "Date"
	}

	if errBuff != "" {
		return "", fmt.Errorf("time entry is missing required field(s): %s", errBuff)
	}

	timeEntryJSON := TimeEntryJSON{
		Entry: entry,
	}

	endpoint := "tasks/" + entry.TaskID + "/time_entries"

	var status struct {
		Status string `json:"STATUS"`
		ID     string `json:"timeLogId"`
	}

	data, err := json.Marshal(timeEntryJSON)
	if err != nil {
		return "", err
	}

	res, err := conn.PostRequest(endpoint, data)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(res, &status)
	if err != nil {
		return "", err
	}

	if status.Status != "OK" {
		return status.Status, fmt.Errorf("received error response (%s)", status.Status)
	}

	if status.ID == "" {
		return "", fmt.Errorf("no id returned for time log")
	}

	entry.ID = status.ID

	return status.ID, nil
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

package teamworkapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
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

// TimeEntriesJSON models the parent JSON structure of an array of TimeEntrys and
// facilitates unmarshalling.
type TimeEntriesJSON struct {
	TimeEntries []*TimeEntry `json:"time-entries"`
}

// TimeQueryParams defines valid query parameters for this resource.
type TimeQueryParams struct {
	UserID 		string `url:"userId,omitempty"`
	FromDate	string `url:"fromdate,omitempty"`
	ToDate		string `url:"todate,omitempty"`
}

// FormatQueryParams formats query parameters for this resource.
func (qp *TimeQueryParams) FormatQueryParams() (string, error) {

	if qp.FromDate != "" {
		_, err := time.Parse("20060102", qp.FromDate)
		if err != nil {
			return "", fmt.Errorf("invalid format for FromDate parameter.  Should be YYYYMMDD, but found %s", qp.FromDate)
		}
	}

	if qp.ToDate != "" {
		_, err := time.Parse("20060102", qp.ToDate)
		if err != nil {
			return "", fmt.Errorf("invalid format for ToDate parameter.  Should be YYYYMMDD, but found %s", qp.ToDate)
		}
	}

	s, err := query.Values(qp)
	if err != nil {
		return "", err
	}

	return s.Encode(), nil
}

// GetTimeEntriesByPerson retrieves time entries for a specific Teamwork user, for the specified time period.
func (conn Connection) GetTimeEntriesByPerson(personID string, fromDate string, toDate string) ([]*TimeEntry, error) {
	
	errBuff := ""

	if personID == "" {
		errBuff += "personID"
	} else {
		_, err := strconv.Atoi(personID)
		if err != nil {
			return nil, err
		}
	}

	if fromDate == "" {
		if errBuff != "" {
			errBuff += ", "
		}
		errBuff += "fromDate"
	}

	if toDate == "" {
		if errBuff != "" {
			errBuff += ", "
		}
		errBuff += "toDate"
	}

	if errBuff != "" {
		return nil, fmt.Errorf("missing required parameter(s): %s", errBuff)
	}

	queryParams := TimeQueryParams {
		UserID: personID,
		FromDate: fromDate,
		ToDate: toDate,
	}

	data, err := conn.GetRequest("time_entries", &queryParams)

	if err != nil {
		return nil, err
	}

	t := new(TimeEntriesJSON)

	err = json.Unmarshal(data, &t)

	if err != nil {
		return nil, err
	}

	return t.TimeEntries, nil
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
func SumHours(e []*TimeEntry, personID string) (float64, error) {
	found := false
	hours := 0.0

	for _, v := range e {
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

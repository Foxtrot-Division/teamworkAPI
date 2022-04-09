package teamworkapi

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

// TeamworkDateFormatShort is the short-form of a date used by Teamwork.
const TeamworkDateFormatShort = "20060102"

// TeamworkDateFormatMed is the medium-form of a date/time used by Teamwork.
const TeamworkDateFormatMed = "2006-01-02T15:04"

// TeamworkDateFormatLong is the long-form of a date/time used by Teamwork.
const TeamworkDateFormatLong = "2006-01-02T15:04:05Z"

// TimeEntry models an individual time entry.
type TimeEntry struct {
	ID          string `json:"id"`
	PersonID    string `json:"person-id"`
	Lastname    string `json:"person-last-name"`
	Firstname   string `json:"person-first-name"`
	Description string `json:"description"`
	Hours       string `json:"hours"`
	Minutes     string `json:"minutes"`
	Date        string `json:"date"` // expected format is YYYYMMDD
	IsBillable  string `json:"isbillable"`
	ProjectID   string `json:"project-id"`
	TaskID      string `json:"todo-item-id"`
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

// TimeResponseHandler models a http response for a TimeEntry operation.
type TimeResponseHandler struct {
	Status      string `json:"STATUS"`
	Message     string `json:"MESSAGE"`
	TimeEntryID string `json:"timeLogId"`
}

type TimeResponseHandlerV3 struct {
	Status  string      `json:"STATUS"`
	Message string      `json:"MESSAGE"`
	TimeLog []TimeLogV3 `json:"timelogs"`
}

// TimeQueryParams defines valid query parameters for this resource.
type TimeQueryParams struct {
	UserID   string `url:"userId,omitempty"`
	FromDate string `url:"fromdate,omitempty"`
	ToDate   string `url:"todate,omitempty"`
	PageSize string `url:"pageSize,omitempty"`
}
type TimeQueryParamsV3 struct {
	EndDate   string `url:"endDate,omitempty"`
	StartDate string `url:"startDate,omitempty"`
	ProjectID string `url:"projectId,omitempty"`
}
type TimeLogV3 struct {
	ID      int `json:"userId"`
	Minutes int `json:"minutes"`
	TaskID  int `json:"taskId"`
}
type TimeLogJSON struct {
	TimeLog []*TimeLogV3 `json:"timelogs"`
}

// ParseResponse interprets a http response for a TimeEntry operation such as
// POST, PUT, UPDATE
func (resMsg *TimeResponseHandler) ParseResponse(httpMethod string, rawRes []byte) error {

	err := json.Unmarshal(rawRes, &resMsg)
	if err != nil {
		return err
	}

	if resMsg.Status == "Error" {
		return fmt.Errorf("received ERROR response: %s", resMsg.Message)
	}

	switch httpMethod {
	case http.MethodPost:
		if resMsg.TimeEntryID == "" {
			return fmt.Errorf("no ID returned for time entry POST")
		}
	}

	return nil
}

// func (resMsg *TimeResponseHandlerV3) ParseResponseV3(httpMethod string, rawRes []byte) error {

// 	b := string(rawRes)
// 	fmt.Println(b)

// 	fmt.Println("raw resp")
// 	err := json.Unmarshal(rawRes, &resMsg)
// 	if err != nil {
// 		return err
// 	}

// 	if resMsg.Status == "Error" {
// 		return fmt.Errorf("received ERROR response: %s", resMsg.Message)
// 	}

// 	switch httpMethod {
// 	case http.MethodPost:
// 		if resMsg.TimeLog[0].ID == 0 {
// 			return fmt.Errorf("no time entry returned")
// 		}
// 	}

// 	//fmt.Println(resMsg.TimeLog[0].Minutes)

// 	return nil
// }

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
func (qp *TimeQueryParamsV3) FormatQueryParamsV3() (string, error) {

	if qp.EndDate != "" {
		_, err := time.Parse("2006-01-02", qp.EndDate)
		if err != nil {
			return "", fmt.Errorf("invalid format for FromDate parameter.  Should be YYYYMMDD, but found %s", qp.EndDate)
		}
	}
	s, err := query.Values(qp)
	if err != nil {
		return "", err
	}

	return s.Encode(), nil
}

// GetTimeEntries retrieve time entries specified by queryParams.
func (conn *Connection) GetTimeEntries(queryParams *TimeQueryParams) ([]*TimeEntry, error) {

	data, err := conn.GetRequest("time_entries", queryParams)
	if err != nil {
		return nil, err
	}
	entries := new(TimeEntriesJSON)
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}
	return entries.TimeEntries, nil
}

func (conn *Connection) GetTimeEntriesV3(queryParams *TimeQueryParamsV3) ([]*TimeLogV3, error) {

	data, err := conn.GetRequestV3("time", queryParams)
	if err != nil {
		return nil, err
	}
	entries := new(TimeLogJSON)
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}
	return entries.TimeLog, nil
}

// GetTimeEntriesByTask retrieves all time entries for the specified Task.
func (conn *Connection) GetTimeEntriesByTask(ID string) ([]*TimeEntry, error) {

	_, err := strconv.Atoi(ID)
	if err != nil {
		if ID == "" {
			return nil, fmt.Errorf("missing required parameter(s): ID")
		}
		return nil, fmt.Errorf("invalid value (%s) for ID", ID)
	}

	endpoint := fmt.Sprintf("tasks/%s/time_entries", ID)

	data, err := conn.GetRequest(endpoint, nil)

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

	queryParams := TimeQueryParams{
		UserID:   personID,
		FromDate: fromDate,
		ToDate:   toDate,
	}

	return conn.GetTimeEntries(&queryParams)
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

	data, err := json.Marshal(timeEntryJSON)
	if err != nil {
		return "", err
	}

	handler := new(TimeResponseHandler)

	err = conn.PostRequest(endpoint, data, handler)
	if err != nil {
		return "", err
	}

	entry.ID = handler.TimeEntryID

	return handler.TimeEntryID, nil
}

// DeleteTimeEntry deletes a time entry with the specified ID.
func (conn *Connection) DeleteTimeEntry(ID string) error {

	if ID == "" {
		return fmt.Errorf("missing required parameter: ID")
	}

	endpoint := "time_entries/" + ID

	handler := new(TimeResponseHandler)

	err := conn.DeleteRequest(endpoint, handler)
	if err != nil {
		return err
	}

	return nil
}

// TotalAndAvgHours returns the total and avg hours found in the TimeEntries array.
func TotalAndAvgHours(e []*TimeEntry) (map[string]float64, error) {
	hours := 0.0

	retVal := make(map[string]float64, 2)
	retVal["total"] = 0.0
	retVal["avg"] = 0.0

	for _, v := range e {
		h, err := strconv.ParseFloat(v.Hours, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to float64", v.Hours)
		}

		hours += h

		if v.Minutes != "0" {
			m, err := strconv.ParseFloat(v.Minutes, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to float64", v.Minutes)
			}
			hours += m / 60
		}
	}

	retVal["total"] = math.Round(hours*100) / 100
	retVal["avg"] = math.Round(hours/float64(len(e))*100) / 100

	return retVal, nil
}

// DurationInDays gets the number of days between the specified dates.
func DurationInDays(from string, to string) (int, error) {

	start, err := time.Parse(TeamworkDateFormatShort, from)
	if err != nil {
		return 0, err
	}

	end, err := time.Parse(TeamworkDateFormatShort, to)
	if err != nil {
		return 0, err
	}

	return int(end.Sub(start).Hours() / 24), nil
}

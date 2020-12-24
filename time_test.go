package teamworkapi

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

type TimeTestData struct {
	People      []string            `json:"people"`
	TimePeriods []map[string]string `json:"time-periods"`
}

func initTimeTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func initTimeTestData(t *testing.T) *TimeTestData {
	testData := new(TimeTestData)

	f, err := os.Open("./testdata/timeTestData.json")
	defer f.Close()

	if err != nil {
		t.Fatalf(err.Error())
	}

	raw, _ := ioutil.ReadAll(f)

	err = json.Unmarshal(raw, &testData)
	if err != nil {
		t.Fatalf(err.Error())
	}

	return testData
}

func TestGetTimeEntriesByPerson(t *testing.T) {

	conn := initTimeTestConnection(t)

	testData := initTimeTestData(t)

	testDateLayout := "20060102"

	for _, p := range testData.People {
		for _, tp := range testData.TimePeriods {

			fromDate, err := time.Parse(testDateLayout, tp["fromdate"])
			if err != nil {
				t.Errorf(err.Error())
			}

			toDate, err := time.Parse(testDateLayout, tp["todate"])
			if err != nil {
				t.Errorf(err.Error())
			}

			entries, err := conn.GetTimeEntriesByPerson(p, tp["fromdate"], tp["todate"])

			if err != nil {
				t.Errorf(err.Error())
			}

			if len(entries) < 1 {
				t.Errorf("No time entries for person %s, from %s to %s.", p, tp["fromdate"], tp["todate"])
			}

			for _, entry := range entries {

				if entry.PersonID != p {
					t.Errorf("Found user ID (%s) but expected only (%s)", entry.PersonID, p)
				}

				entryTime, err := time.Parse(time.RFC3339, entry.Date)
				if err != nil {
					t.Errorf(err.Error())
				}
				d := time.Date(entryTime.Year(), entryTime.Month(), entryTime.Day(), 0, 0, 0, 0, time.UTC)

				if d.Before(fromDate) || d.After(toDate) {
					t.Errorf("Entry (%s) is not within specified time range (%s - %s)!", d, fromDate, toDate)
				}
			}
		}
	}

	var tests = []struct {
		ID 		string
		from 	string
		to 		string
		want	string
	}{
		{"", "20201012", "20201013", "missing required parameter(s): personID"},
		{"abc", "", "", "strconv.Atoi: parsing \"abc\": invalid syntax"},
		{"12345", "10-12-2020", "20201013", "invalid format for FromDate parameter.  Should be YYYYMMDD, but found 10-12-2020"},
	}

	for _, v := range tests {
		_, err := conn.GetTimeEntriesByPerson(v.ID, v.from, v.to)
		if err != nil {
			if err.Error() != v.want {
				t.Errorf("expected error (%s) but got (%s)", v.want, err.Error())
			}
		} else {
			t.Errorf("Expected error for userID (%s)", v.ID)
		}
	}		
}

func TestPostTimeEntry(t *testing.T) {
	
	conn := initTimeTestConnection(t)

	var tests = []struct {
		personID 	string
		description string
		hours 		string
		minutes		string
		date		string
		isBillable  string
		taskID 		string
		error		bool
		want		string
	}{
		{"118616", "Test entry 2.", "10", "", "20201222", "false", "20029437", false, ""},
		{"", "Test entry.", "0", "0", "", "true", "", true, "time entry is missing required field(s): PersonID, TaskID, Date"},
	}

	for _, v := range tests {

		entry := &TimeEntry{
			PersonID: v.personID,
			Description: v.description,
			Hours: v.hours,
			Minutes: v.minutes,
			Date: v.date,
			IsBillable: v.isBillable,
			TaskID: v.taskID,
		}
		
		res, err := conn.PostTimeEntry(entry)

		if err != nil {
			if !v.error {
				t.Errorf(err.Error())
			} else {
				if err.Error() != v.want {
					t.Errorf("expected error (%s) but got (%s)", v.want, err.Error())
				}
			}
		} else {
			if v.error {
				t.Errorf("expected error")
			} else {
				if entry.ID != res {
					t.Errorf("ID (%s) not set to expected value (%s)", entry.ID, res)
				}
			}
		}
	}		
}

func TestSumHours(t *testing.T) {

	conn := initTimeTestConnection(t)
	testData := initTimeTestData(t)

	for _, p := range testData.People {
		for _, tp := range testData.TimePeriods {

			entries, err := conn.GetTimeEntriesByPerson(p, tp["fromdate"], tp["todate"])

			if err != nil {
				t.Errorf(err.Error())
			}

			hours, err := SumHours(entries, p)

			if err != nil {
				t.Errorf(err.Error())
			}

			if hours < 1 {
				t.Errorf("No hours found for user ID %s", p)
			}
		}
	}
}

package teamworkapi

import (
	"encoding/json"
	"fmt"
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
	conn, err := NewConnectionFromJSON("./testdata/apiConfig.json")
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

			if len(entries.TimeEntries) < 1 {
				t.Errorf("No time entries for person %s, from %s to %s.", p, tp["fromdate"], tp["todate"])
			}

			for _, entry := range entries.TimeEntries {

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

	_, err := conn.GetTimeEntriesByPerson("1234", "badformat", "2010-02-03")
	if err == nil {
		t.Errorf("invalid string allowed for from/to parameter")
	}
}

func TestPostTimeEntry(t *testing.T) {
	
	conn := initTimeTestConnection(t)

	entry := TimeEntry{
			PersonID: "118616",
			Description: "Test entry.",
			Hours: "0",
			Minutes: "60",
			Date: "20201209",
			IsBillable: "false",
		}

	res, err := conn.PostTimeEntry("20029437", entry)

	if err != nil {
		t.Errorf(err.Error())
	}

	entry.ID = res

	fmt.Printf("TimeEntry ID: %s", entry.ID)
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

			hours, err := entries.SumHours(p)

			if err != nil {
				t.Errorf(err.Error())
			}

			if hours < 1 {
				t.Errorf("No hours found for user ID %s", p)
			}
		}
	}
}

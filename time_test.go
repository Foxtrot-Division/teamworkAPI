package teamworkapi

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

type TimeTestData struct {
	People 		[]string `json:"people"`
	TimePeriods	[] map[string] string `json:"time-periods"`
}

func initTestConnection(t *testing.T) *connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfig.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func initTestData(t *testing.T) *TimeTestData {
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
	conn := initTestConnection(t)

	testData := initTestData(t)

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

				entryTime := entry.Date
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

func TestSumHours(t *testing.T) {

	conn := initTestConnection(t)
	testData := initTestData(t)

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
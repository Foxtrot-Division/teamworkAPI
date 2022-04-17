package teamworkapi

import (
	"fmt"
	"testing"
	"os"
	"io/ioutil"
	"encoding/json"
)

func initCalendarEventTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func initCalendarConnectionV3(t *testing.T) *Connection {

	f, err := os.Open("./testdata/tw_api_conf.json")
	defer f.Close()

	if err != nil {
		t.Errorf(err.Error())
	}

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	data := new(TWAPIConf)

	err = json.Unmarshal(raw, &data)
	if err != nil {
		t.Errorf(err.Error())
	}

	//	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	conn, err := NewConnection(data.APIKey, data.SiteName, "", data.APIVersion)
	if err != nil {
		t.Fatalf(err.Error())
	}
	
	return conn
}

func TestGetCalendarEvents(t *testing.T) {

	conn := initCalendarEventTestConnection(t)

	q := CalendarEventQueryParams{
		From: "20210101",
	}

	events, err := conn.GetCalendarEvents(q)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(events) < 1 {
		t.Errorf("no events retrieved")
	}
}

func TestGetCalendarEventsV3(t *testing.T) {

	conn := initCalendarConnectionV3(t)

	q := CalendarEventQueryParamsV3{
		StartDate: "2021-01-01",
		EndDate:   "2022-01-01",
		PageSize: "1000",
	}

	events, err := conn.GetCalendarEventsV3(q)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(events) < 1 {
		t.Errorf("no events retrieved")
	}

	for _, u := range events {
		fmt.Printf("Start Date: %v, ", *&u.StartDate)
		fmt.Printf("End Date: %v, ", *&u.EndDate)
		fmt.Printf("Type ID: %v, ", *&u.TypeId)
		fmt.Printf("Owner User ID: %v, ", *&u.OwnerUserId)
		fmt.Printf("All Day: %v, ", *&u.AllDay)
		fmt.Printf("Attending User ID: %v \n", *&u.AttendingUserIds)

	}
}

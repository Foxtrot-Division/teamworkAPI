package teamworkapi

import (
	"fmt"
	"testing"
)

func initCalendarEventTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func initCalenderTestConnection(t *testing.T) *Connection {
	conn, err := NewConnection("water589meat", "foxtrotdivision", "", "v3")
	if err != nil {
		t.Fatalf(err.Error())
	}
	//fmt.Print("Connection")
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

	conn := initCalenderTestConnection(t)

	q := CalendarEventQueryParamsV3{
		StartDate: "2021-01-01",
		EndDate:   "2022-01-01",
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
		fmt.Printf("Attending User ID: %v \n", *&u.AttendingUserIds)

	}
}

package teamworkapi

import (

	"testing"
)



func initCalendarEventTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func TestGetCalendarEvents(t *testing.T) {

	conn := initCalendarEventTestConnection(t)

	q := CalendarEventQueryParams {
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
package teamworkapi

import (
	"testing")


func initProjectTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func TestGetProjects(t *testing.T) {

	conn := initProjectTestConnection(t)

	projects, err := conn.GetProjects(nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	want := 14

	if len(projects) != want {
		t.Errorf("expected %d projects but got %d", want, len(projects))
	}
}
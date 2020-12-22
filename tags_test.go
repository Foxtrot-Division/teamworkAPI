package teamworkapi

import (
	"testing"
)

func initTagsTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}
func TestGetTags(t *testing.T) {

	conn := initTagsTestConnection(t)

	tags, err := conn.GetTags()
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(tags) < 1 {
		t.Errorf("no tags returned")
	}
}
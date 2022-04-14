package teamworkapi

import (
	"testing"
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


func initProjectTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func initProjectTestConnectionV3(t *testing.T) *Connection {
	
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

func TestGetProjectV3(t *testing.T) {

	conn := initProjectTestConnectionV3(t)
	
	project, err := conn.GetProjectV3("223404")
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(project)
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
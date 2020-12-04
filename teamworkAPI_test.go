package teamworkapi

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewConnection(t *testing.T) {
	key := "someKey"
	site := "someSite"
	dataPref := "json"

	conn, err := NewConnection(key, site, dataPref)

	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}

	if conn.APIKey != key {
		t.Errorf("apiKey (%s) should be (%s)", conn.APIKey, key)
	}

	if conn.SiteName != site {
		t.Errorf("siteName (%s) should be (%s)", conn.SiteName, site)
	}

	if conn.DataPreference != dataPref {
		t.Errorf("dataPreference (%s) should be (%s)", conn.DataPreference, dataPref)
	}

	conn, err = NewConnection(key, site, "")

	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}

	if conn.DataPreference != "json" {
		t.Errorf("dataPreference (%s) should have defaulted to json", conn.DataPreference)
	}

	conn, err = NewConnection("", "", "")

	if err == nil {
		t.Errorf("NewTeamworkAPI call was allowed with empty string parameter values.")
	}
}

func TestNewTeamworkAPIFromJSON(t *testing.T) {
	conn, err := NewConnectionFromJSON("./testdata/apiConfig.json")

	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}

	if conn.URL != "https://" + conn.SiteName + ".teamwork.com/" {
		t.Errorf("URL (%s) not formed correctly", conn.URL)
	}
}

func TestGetRequest(t *testing.T) {
	var res struct {
		Status string `json:"STATUS"`
	}

	conn, _ := NewConnectionFromJSON("./testdata/apiConfig.json")

	data, err := conn.GetRequest("projects")

	if err != nil {
		t.Errorf(fmt.Sprintf("%s", err))
	}

	err = json.Unmarshal(data, &res)

	if err != nil {
		t.Errorf(err.Error())
	}

	if res.Status != "OK" {
		t.Errorf("Received STATUS (%s) when it should be (OK)", res.Status)
	}
}

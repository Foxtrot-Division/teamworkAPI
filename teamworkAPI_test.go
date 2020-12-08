package teamworkapi

import (
	"encoding/json"
	"fmt"
	"net/url"
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

	if conn.URL != "https://"+conn.SiteName+".teamwork.com/" {
		t.Errorf("URL (%s) not formed correctly", conn.URL)
	}
}

func TestGetRequest(t *testing.T) {

	var raw interface{}

	var tests = []struct {
		endpoint string
		params   map[string]interface{}
		expect   string
	}{
		{"projects", nil, "OK"},
		{"people", nil, "OK"},
		{"companies", nil, "OK"},
	}

	conn, _ := NewConnectionFromJSON("./testdata/apiConfig.json")

	for _, tt := range tests {

		data, err := conn.GetRequest(tt.endpoint, tt.params)

		if err != nil {
			t.Errorf(fmt.Sprintf("%s", err))
		}

		err = json.Unmarshal(data, &raw)

		if err != nil {
			t.Errorf(err.Error())
		}

		res := raw.(map[string]interface{})

		if res["STATUS"] != tt.expect {
			t.Errorf("Received STATUS (%s) but expected (%s)", res["STATUS"], tt.expect)
		}

	}
}

func TestFormatQueryString(t *testing.T) {

	var tests = []struct {
		params map[string]interface{}
		expect string
	}{
		{map[string]interface{}{"key1": "val1", "key2": true, "key3": 10}, url.Values{"key1": []string{"val1"}, "key2": []string{"true"}, "key3": []string{"10"}}.Encode()},
		{map[string]interface{}{"key1": false, "key2": "val-2", "key3": 133}, url.Values{"key1": []string{"false"}, "key2": []string{"val-2"}, "key3": []string{"133"}}.Encode()},
		{map[string]interface{}{"key1": "#val3", "key2": true, "key3": 0}, url.Values{"key1": []string{"#val3"}, "key2": []string{"true"}, "key3": []string{"0"}}.Encode()},
	}

	for _, tt := range tests {
		result, err := FormatQueryString(tt.params)

		if err != nil {
			t.Errorf(err.Error())
		}

		if result.Encode() != tt.expect {
			t.Errorf("Expected: %s\nGot: %s\n", result.Encode(), tt.expect)
		}

	}

}

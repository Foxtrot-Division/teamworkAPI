package teamworkapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func TestNewConnection(t *testing.T) {

	// test error conditions
	var tests1 = []struct {
		key			string
		site 		string
		dataPref 	string
		err  		bool
		want 		string
	}{
		{key: "someKey", site: "someSite", dataPref: "json", err: false, want: ""},
		{key: "123456!$", site:"", dataPref:"json", err: true, want: "missing required parameter(s): siteName"},
		{key: "", site:"", dataPref:"", err: true, want:"missing required parameter(s): apiKey, siteName"},
		{key: "buddha", site:"belly", dataPref:"", err: false, want:""},
	}

	// test default setting for dataPreference
	var tests2 = []struct {
		key			string
		site 		string
		dataPref 	string
		want 		string
	}{
		{key: "someKey", site: "someSite", dataPref: "", want: "json"},
		{key: "vader", site:"deathstar", dataPref:"json", want: "json"},
		{key: "gold", site:"bravo", dataPref:"someFormat", want: "someFormat"},
	}

	for _, v := range tests1 {

		conn, err := NewConnection(v.key, v.site, v.dataPref)

		if err != nil {
			if !v.err {
				t.Errorf(err.Error())
			} else {
				if v.want != err.Error() {
					t.Errorf("expected error string (%s) but got (%s)", v.want, err.Error())
				}
			}
		} else {
			if v.err {
				t.Errorf("expected error for input (key: %s, site: %s, dataPref: %s)", v.key, v.site, v.dataPref)
			} else {
				if conn.APIKey != v.key {
					t.Errorf("expected APIKey (%s) but got (%s)", v.key, conn.APIKey)
				}
		
				if conn.SiteName != v.site {
					t.Errorf("expected SiteName(%s) but got (%s)", v.site, conn.SiteName)
				}
			}
		}
	}

	for _, v := range tests2 {

		conn, err := NewConnection(v.key, v.site, v.dataPref)

		if err != nil {
			t.Errorf(err.Error())
		}

		if conn.DataPreference != v.want {
			t.Errorf("expected DataPreference (%s) but got (%s)", v.want, conn.DataPreference)
		}
	}
}

func TestNewTeamworkAPIFromJSON(t *testing.T) {

	// test error conditions
	var tests1 = []struct {
		fileName	string
		err  		bool
		want 		string
	}{
		{fileName: "apiConfigTestData1.json", err: false, want: ""},
		{fileName: "apiConfigTestData2.json", err: true, want: "missing required parameter(s): apiKey"},
		{fileName: "apiConfigTestData3.json", err: true, want:"missing required parameter(s): apiKey, siteName"},
		{fileName: "badFileName.json", err: true, want:"Failed to open JSON file at ./testdata/badFileName.json"},
	}

	// test default setting for dataPreference
	var tests2 = []struct {
		fileName	string
		want 		string
	}{
		{fileName: "apiConfigTestData4.json", want: "json"},
		{fileName: "apiConfigTestData5.json", want: "someFormat"},
	}

	for _, v := range tests1 {

		conn, err := NewConnectionFromJSON("./testdata/" + v.fileName)

		if err != nil {
			if !v.err {
				t.Errorf(err.Error())
			} else {
				if v.want != err.Error() {
					t.Errorf("expected error string (%s) but got (%s)", v.want, err.Error())
				}
			}
		} else {
			if v.err {
				t.Errorf("expected error for input (key: %s, site: %s, dataPref: %s)", conn.APIKey, conn.SiteName, conn.DataPreference)
			} else {
				if conn.URL != "https://" + conn.SiteName + ".teamwork.com/" {
					t.Errorf("URL (%s) not formed correctly", conn.URL)
				}
			}
		}
	}

	for _, v := range tests2 {

		conn, err := NewConnectionFromJSON("./testdata/" + v.fileName)

		if err != nil {
			t.Errorf(err.Error())
		}

		if conn.DataPreference != v.want {
			t.Errorf("expected DataPreference (%s) but got (%s)", v.want, conn.DataPreference)
		}
	}
}

func TestGetRequest(t *testing.T) {

	var raw interface{}

	// test sample of good/bad endpoints
	var tests = []struct {
		endpoint 	string
		params   	map[string]interface{}
		want   		string
	}{
		{"projects", nil, "OK"},
		{"people", nil, "OK"},
		{"companies", nil, "OK"},
		{"buffalo", nil, ""},
		{"people", map[string]interface{}{"sort":"company"}, "OK"},
		{"projects", map[string]interface{}{"status":"ACTIVE", "includePeople": true}, "OK"},
		{"tasks", map[string]interface{}{"sort":"project", "includeArchivedProjects": true}, "OK"},
	}

	conn, _ := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")

	for _, tt := range tests {

		data, err := conn.GetRequest(tt.endpoint, tt.params)

		if err != nil {
			t.Errorf(fmt.Sprintf("%s", err))
		}

		err = json.Unmarshal(data, &raw)

		if err != nil {
			fmt.Println(string(data))
			t.Errorf(err.Error())
		}

		res := raw.(map[string]interface{})

		if res["STATUS"] != tt.want {
			if res["STATUS"] == nil && tt.want != "" {
				t.Errorf("Received response (%s) but expected (%s)", res["STATUS"], tt.want)
			}
		}

	}
}

func TestFormatQueryString(t *testing.T) {

	var tests = []struct {
		params 	map[string]interface{}
		want 	string
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

		if result.Encode() != tt.want {
			t.Errorf("Expected: %s\nGot: %s\n", result.Encode(), tt.want)
		}
	}
}

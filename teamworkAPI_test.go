package teamworkapi

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-querystring/query"
)

type GenericQueryParams struct {
	Sort string `url:"sort,omitempty"`
	Status string `url:"status,omitempty"`
	IncludePeople bool `url:"includePeople,omitempty"`
	IncludeArchivedProjects bool `url:"includeArchivedProjects,omitempty"`
}

func (params GenericQueryParams) FormatQueryParams() (string, error) {
	
	p, err := query.Values(params)
	if err != nil {
		return "", err
	}

	return p.Encode(), nil
}

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
		params   	GenericQueryParams
		want   		string
	}{
		{"projects", GenericQueryParams{}, "OK"},
		{"people", GenericQueryParams{}, "OK"},
		{"companies", GenericQueryParams{}, "OK"},
		{"buffalo", GenericQueryParams{}, ""},
		{"people", GenericQueryParams {Sort: "company"}, "OK"},
		{"projects", GenericQueryParams {Status: "ACTIVE", IncludePeople: true}, "OK"},
		{"tasks", GenericQueryParams{Sort: "project", IncludeArchivedProjects: true}, "OK"},
	}

	conn, _ := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")

	for i, tt := range tests {

		data, err := conn.GetRequest(tt.endpoint, tt.params)

		if err != nil {
			t.Errorf(fmt.Sprintf("%s", err))
		}

		err = json.Unmarshal(data, &raw)

		if err != nil {
			t.Errorf(err.Error())
		}

		res := raw.(map[string]interface{})

		if res["STATUS"] != tt.want {
			if res["STATUS"] == nil && tt.want != "" {
				t.Errorf("test [%d] received response (%s) but expected (%s)", i, res["STATUS"], tt.want)
			}
		}
	}
}

package teamworkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type peopleTestData struct {
	People []*Person	`json:"people"`
}

func initPeopleTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func loadPeopleTestData(t *testing.T) *peopleTestData {

	f, err := os.Open("./testdata/peopleTestData.json")
	defer f.Close()
	
	if err != nil {
		t.Errorf(err.Error())
	}

	data := new(peopleTestData)
	
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		t.Errorf(err.Error())
	}

	return data
}

func TestGetPersonByID(t *testing.T) {
	
	testData := loadPeopleTestData(t)

	conn := initPeopleTestConnection(t)

	// test valid cases
	for _, v := range testData.People {

		p, err := conn.GetPersonByID(v.ID)
		if err != nil {
			t.Errorf(err.Error())
		}

		if p == nil {
			t.Errorf("No data returned for ID (%s)", v.ID)
		} else {
			if p.Email != v.Email {
				t.Errorf("Expected email (%s) but got (%s)", v.Email, p.Email)
			}
			
			if p.FirstName != v.FirstName {
				t.Errorf("Expected FirstName (%s) but got (%s)", v.FirstName, p.FirstName)
			}

			if p.LastName != v.LastName {
				t.Errorf("Expected LastName (%s) but got (%s)", v.LastName, p.LastName)
			}

			if p.CompanyName != v.CompanyName {
				t.Errorf("Expected CompanyName (%s) but got (%s)", v.CompanyName, p.CompanyName)
			}
		}
	}

	// test error responses
	var tests = []struct {
		ID 		string
		want	string
	}{
		{"123456", "failed to retrieve user with ID (123456)"},
		{"bad-content", "invalid value (bad-content) for ID"},
		{"", "missing required parameter(s): ID"},
	}

	for _, v := range tests {
		_, err := conn.GetPersonByID(v.ID)
		if err != nil {
			if err.Error() != v.want {
				t.Errorf("expected error (%s) but got (%s)", v.want, err.Error())
			}
		} else {
			t.Errorf("Expected error for userID (%s)", v.ID)
		}
	}		
}

func TestGetPeople(t *testing.T) {
	conn, err := NewConnectionFromJSON("./testdata/apiConfig.json")

	if err != nil {
		t.Errorf(err.Error())
	}

	var c map[string] interface{}

	f, err := os.Open("./testdata/companyTestData.json")
	defer f.Close()
	
	if err != nil {
		t.Errorf(err.Error())
	}
	
	data, _ := ioutil.ReadAll(f)
	
	err = json.Unmarshal(data, &c)
	if err != nil {
		t.Fatal(err.Error())
	}

	id := fmt.Sprintf("%v", c["company-with-people"])
	
	p, err := conn.GetPeopleByCompany(string(id))

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(p.People) < 1 {
		t.Errorf("No people returned for company ID %s", id)
	}
}

func TestGetCompanies(t *testing.T) {

	conn := initPeopleTestConnection(t)

	c, err := conn.GetCompanies()

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(c.Companies) < 1 {
		t.Errorf("No companies returned.")
	}
}
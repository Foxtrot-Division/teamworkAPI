package teamworkapi

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

type peopleTestData struct {
	ExampleUserID 		string 	`json:"exampleUserID"`
	CompanyID   		string 	`json:"companyID"`
	CompanyTotalUsers	int 	`json:"companyTotalUsers"`
	ExampleProjectID 	string 	`json:"exampleProjectID"`
	ProjectTotalUsers 	int		`json:"projectTotalUsers"`
}

type peopleTestDataJSON struct {
	Data []peopleTestData `json:"data"`
}

func initPeopleTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func loadPeopleTestData(t *testing.T) []peopleTestData {

	f, err := os.Open("./testdata/peopleTestData.json")
	defer f.Close()
	
	if err != nil {
		t.Errorf(err.Error())
	}

	data := new(peopleTestDataJSON)
	
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		t.Errorf(err.Error())
	}

	return data.Data
}

func TestGetPersonByID(t *testing.T) {
	
	testData := loadPeopleTestData(t)

	conn := initPeopleTestConnection(t)

	// test valid cases
	for _, v := range testData {

		p, err := conn.GetPersonByID(v.ExampleUserID)
		if err != nil {
			t.Errorf(err.Error())
		}

		if p == nil {
			t.Errorf("No data returned for ID (%s)", v.ExampleUserID)
		} else {
			if p.CompanyName == "" {
				t.Errorf("no company name returned for user ID (%s)", v.ExampleUserID)
			}
			
			if p.Email == "" {
				t.Errorf("no email returned for user ID (%s)", v.ExampleUserID)
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

func TestGetPeopleByCompany(t *testing.T) {

	testData := loadPeopleTestData(t)

	conn := initPeopleTestConnection(t)

	for _, v := range testData {

		p, err := conn.GetPeopleByCompany(v.CompanyID)

		if err != nil {
			t.Errorf(err.Error())
		}
	
		if len(p) != v.CompanyTotalUsers {
			t.Errorf("expected (%d) users for company ID (%s) but got (%d)", v.CompanyTotalUsers, v.CompanyID, len(p))
		}
	}
}

func TestGetPeople(t *testing.T) {

	testData := loadPeopleTestData(t)

	conn := initPeopleTestConnection(t)

	userIDs := ""
	numberUsers := 0

	for _, v := range testData {

		if userIDs != "" {
			userIDs += ","
		}

		userIDs += v.ExampleUserID
		
		q1 := PeopleQueryParams {
			ProjectID: v.ExampleProjectID,
		}
	
		people, err := conn.GetPeople(q1)
		if err != nil {
			t.Errorf(err.Error())
		}
	
		if len(people) != v.ProjectTotalUsers {
			t.Errorf("expected (%d) users but got (%d) %s", v.ProjectTotalUsers, len(people), conn.RequestURL)
		}

		numberUsers++
	}

	q2 := PeopleQueryParams {
		UserID: userIDs,
	}

	people, err := conn.GetPeople(q2)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(people) != numberUsers {
		t.Errorf("expected (%d) users but got (%d) %s", numberUsers, len(people), conn.RequestURL)
	}
}

func TestGetCompanies(t *testing.T) {

	conn := initPeopleTestConnection(t)

	c, err := conn.GetCompanies()

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(c) < 1 {
		t.Errorf("No companies returned.")
	}
}
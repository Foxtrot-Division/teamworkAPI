package teamworkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

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
	
	p, err := conn.GetPeople(string(id))

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(p.People) < 1 {
		t.Errorf("No people returned for company ID %s", id)
	}
}

func TestGetCompanies(t *testing.T) {
	conn, err := NewConnectionFromJSON("./testdata/apiConfig.json")

	if err != nil {
		t.Errorf(err.Error())
	}

	c, err := conn.GetCompanies()

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(c.Companies) < 1 {
		t.Errorf("No companies returned.")
	}
}
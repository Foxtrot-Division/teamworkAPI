package teamworkapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Person models an individual Teamwork user.
type Person struct {
	ID          string `json:"id"`
	FirstName   string `json:"first-name"`
	LastName    string `json:"last-name"`
	CompanyName string `json:"company-name"`
	Email		string `json:"user-name"`
}

// PersonJSON is a wrapper to facilitate marshalling of Person data to json.
type PersonJSON struct {
	Person *Person `json:"person"`
}

// People models an array of individual users.
type People struct {
	People []Person `json:"people"`
}

// Company models an individual company on Teamwork.
type Company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Companies models an array of individual companies.
type Companies struct {
	Companies []Company `json:"companies"`
}

// GetPeopleByCompany retrieves all people from the company specified by companyID.  If companyID is empty string, all people will be returned.
func (conn Connection) GetPeopleByCompany(companyID string) (*People, error) {
	var endpoint = ""

	if companyID != "" {
		endpoint = "companies/" + companyID + "/people"
	} else {
		endpoint = "people"
	}

	data, err := conn.GetRequest(endpoint, nil)

	if err != nil {
		return nil, err
	}

	p := new(People)

	err = json.Unmarshal(data, &p)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetPersonByID retrieves a specific person based on ID. 
func (conn Connection) GetPersonByID(ID string) (*Person, error) {

	_, err := strconv.Atoi(ID)
	if err != nil {
		if ID == "" {
			return nil, fmt.Errorf("missing required parameter(s): ID")
		}
		return nil, fmt.Errorf("invalid value (%s) for ID", ID)
	}

	endpoint := "people/" + ID

	data, err := conn.GetRequest(endpoint, nil)

	if err != nil {
		return nil, err
	}

	p := new(PersonJSON)

	err = json.Unmarshal(data, &p)

	if err != nil {
		return nil, err
	}

	if p.Person == nil {
		return nil, fmt.Errorf("failed to retrieve user with ID (%s)", ID)
	}

	return p.Person, nil
}

// GetCompanies retrieves all companies from Teamwork.
func (conn Connection) GetCompanies() (*Companies, error) {

	data, err := conn.GetRequest("companies", nil)

	if err != nil {
		return nil, err
	}

	c := new(Companies)

	err = json.Unmarshal(data, &c)

	if err != nil {
		return nil, err
	}

	return c, nil
}

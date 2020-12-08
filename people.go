package teamworkapi

import (
	"encoding/json"
)

// Person models an individual Teamwork user.
type Person struct {
	ID          string `json:"id"`
	FirstName   string `json:"first-name"`
	LastName    string `json:"last-name"`
	CompanyName string `json:"company-name"`
}

// People models an array of individual users.
type People struct {
	People []Person `json:"people"`
}

// Company models an individual company on Teamwork.
type Company struct {
	ID		string `json:"id"`
	Name	string `json:"name"`
}

// Companies models an array of individual companies.
type Companies struct {
	Companies []Company `json:"companies"`
}

// GetPeople retrieves all people from the company specified by companyID.  If companyID is empty string, all people will be returned.
func (conn connection) GetPeople(companyID string) (*People, error) {
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

// GetCompanies retrieves all companies from Teamwork.
func (conn connection) GetCompanies() (*Companies, error) {
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
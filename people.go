package teamworkapi

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/go-querystring/query"
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

// PeopleJSON models the parent JSON structure of an array of Persons and
// facilitates unmarshalling.
type PeopleJSON struct {
	People []*Person `json:"people"`
}

// PeopleQueryParams defines valid query parameters for this resource.
type PeopleQueryParams struct {
	ProjectID 	 	 string `url:"projectId,omitempty"`
	UserID	   	 	 string `url:"userIds,omitempty"`
	CompanyID  		 string `url:"companyId,omitempty"`
}

// Company models an individual company on Teamwork.
type Company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CompaniesJSON models the parent JSON structure of an array of Companys and
// facilitates unmarshalling.
type CompaniesJSON struct {
	Companies []*Company `json:"companies"`
}

// FormatQueryParams formats query parameters for this resource.
func (qp PeopleQueryParams) FormatQueryParams() (string, error) {

	params, err := query.Values(qp)
	if err != nil {
		return "", err
	}

	return params.Encode(), nil
}

// GetPeopleByCompany retrieves all people from the company specified by companyID.
func (conn *Connection) GetPeopleByCompany(companyID string) ([]*Person, error) {
	
	_, err := strconv.Atoi(companyID)
	if err != nil {
		if companyID == "" {
			return nil, fmt.Errorf("missing required parameter(s): companyID")
		}
		return nil, fmt.Errorf("invalid value (%s) for companyID", companyID)
	}

	qp := PeopleQueryParams {
		CompanyID: companyID,
	}

	data, err := conn.GetPeople(qp)

	if len(data) < 1 {
		return nil, fmt.Errorf("failed to retrieve any users for companyID (%s)", companyID)
	}

	return data, nil
}

// GetPersonByID retrieves a specific person based on ID. 
func (conn *Connection) GetPersonByID(ID string) (*Person, error) {

	_, err := strconv.Atoi(ID)
	if err != nil {
		if ID == "" {
			return nil, fmt.Errorf("missing required parameter(s): ID")
		}
		return nil, fmt.Errorf("invalid value (%s) for ID", ID)
	}

	qp := PeopleQueryParams {
		UserID: ID,
	}

	data, err := conn.GetPeople(qp)

	if len(data) != 1 {
		return nil, fmt.Errorf("failed to retrieve user with ID (%s)", ID)
	}

	return data[0], nil
}

// GetPeople retrieves people based on query parameters.
func (conn *Connection) GetPeople(queryParams PeopleQueryParams) ([]*Person, error) {
	
	data, err := conn.GetRequest("people", queryParams)
	if err != nil {
		return nil, err
	}

	people := new(PeopleJSON)

	err = json.Unmarshal(data, &people)
	if err != nil {
		return nil, err
	}
	
	return people.People, nil
}

// GetCompanies retrieves all companies from Teamwork.
func (conn *Connection) GetCompanies() ([]*Company, error) {

	data, err := conn.GetRequest("companies", nil)

	if err != nil {
		return nil, err
	}

	c := new(CompaniesJSON)

	err = json.Unmarshal(data, &c)

	if err != nil {
		return nil, err
	}

	return c.Companies, nil
}

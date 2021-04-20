package teamworkapi

import (
	"encoding/json"

	"github.com/google/go-querystring/query"
)

// Project models a Teamwork project.
type Project struct {
	ID          string 	`json:"id"`
	Name    	string 	`json:"name"`
	Description string 	`json:"description"`
	Status 		string 	`json:"status"`
	Company 	Company `json:"company"`
}

// ProjectsJSON provides a wrapper around TimeEntry to properly marshal json
// daProjectsosting to API.
type ProjectsJSON struct {
	Projects []*Project `json:"projects"`
}

// ProjectQueryParams defines valid query parameters for this resource.
type ProjectQueryParams struct {
	CompanyID   string `url:"companyId,omitempty"`
	Status 		string `url:"status,omitempty"`
	PageSize	string `url:"pageSize,omitempty"`
}

// FormatQueryParams formats query parameters for this resource.
func (qp *ProjectQueryParams) FormatQueryParams() (string, error) {

	if qp == nil {
		return "", nil
	}

	s, err := query.Values(qp)
	if err != nil {
		return "", err
	}

	return s.Encode(), nil
}

// GetProjects retrieve projects specified by queryParams.
func (conn *Connection) GetProjects(queryParams *ProjectQueryParams) ([]*Project, error) {

	data, err := conn.GetRequest("projects", queryParams)
	if err != nil {
		return nil, err
	}

	projects := new(ProjectsJSON)
	
	err = json.Unmarshal(data, &projects)
	if err != nil {
		return nil, err
	}

	return projects.Projects, nil
}
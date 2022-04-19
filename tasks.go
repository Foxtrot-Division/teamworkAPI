package teamworkapi

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

// Task models a specific task in Teamwork for Version 3. Refer to TW API docs to add additonal fields as requried.
type TaskVersion3 struct {
	Task struct{
		Id				 int				`json:"id"`
		Description      string             `json:"description"`
		EstimatedMinutes int                `json:"estimatedMinutes"`
		Name             string             `json:"name"`
		Private          bool               `json:"private"`
		ParentTaskID     int                `json:"parentTaskId"`
		Assigness []struct {
			ID int `json:"id"`
			Type string `json:"type"`
		} `json:"assignees"`
		Attachments []struct {
			ID int `json:"id"`
			Type string `json:"type"`
		} `json:"attachments"`
	} `json:"task"`
}

// Task models a specific task in Teamwork.
type Task struct {
	ID             int    `json:"id"`
	Title          string `json:"content"`
	Description    string `json:"description"`
	ProjectID      int    `json:"project-id"`
	TaskListID     int    `json:"todo-list-id"`
	Status         string `json:"status"`
	CompanyID      int    `json:"company-id"`
	DueDate        string `json:"due-date"`
	CreatedOn      string `json:"created-on"`
	CompletedOn    string `json:"completed_on"`
	EstimatedMin   int    `json:"estimated-minutes"`
	Priority       string `json:"priority"`
	AssignedUserID string `json:"responsible-party-id"`
	TimeTotals     *TimeTotals
	Tags           []Tag `json:"tags"`
}

// TaskJSON models the parent JSON structure of an individual task and
// facilitates unmarshalling.
type TaskJSON struct {
	Task *Task `json:"todo-item"`
}

// TasksJSON models the parent JSON structure of an array of tasks and
// facilitates unmarshalling.
type TasksJSON struct {
	Tasks []*Task `json:"todo-items"`
}

type TaskV3 struct {
	Id				 int				`json:"id"`
	Description      string             `json:"description"`
	EstimatedMinutes int                `json:"estimatedMinutes"`
	Name             string             `json:"name"`
	Private          bool               `json:"private"`
	ParentTaskID     int                `json:"parentTaskId"`
	Assignees        map[string][]int64 `json:"assignees"`
}

type TaskV3JSON struct {
	Task TaskV3 `json:"task"`
}

type TaskResponseV3 struct {
	ID int `json:"id"`
}

type TaskPatchV3JSON struct{
	Attachments TaskPatchAttachments `json:"attachments"`
}

type TaskPatchAttachments struct{
	PendingFiles []TaskPatchPendingFiles `json:"pendingFiles"`
}

type TaskPatchPendingFiles struct{
	CategoryId int `json:"categoryId"`
	Reference string `json:"reference"`
}

// TaskResponseHandlerV3 models a http response for a Task operation using version 3 of teamwork api.
type TaskResponseHandlerV3 struct {
	Status  string         `json:"STATUS"`
	Message string         `json:"MESSAGE"`
	Task    TaskResponseV3 `json:"task`
}

// TimeTotals summarizes actual and estimated hours for a specific task.
type TimeTotals struct {
	ActualHours    float64
	EstimatedHours float64
	PercentError   float64
}

// TaskTimeTotalJSON is used to unmarshal the json response provided by call to
// Teamwork API endpoint /tasks/{id}/time/total.json.
type TaskTimeTotalJSON struct {
	Tasklist struct {
		Task struct {
			TimeEstimates struct {
				EstimatedHours string `json:"total-hours-estimated"`
			} `json:"time-estimates"`
			TimeTotals struct {
				ActualHours string `json:"total-hours-sum"`
			} `json:"time-totals"`
		} `json:"task"`
	} `json:"tasklist"`
}

// TaskTimeTotalsJSON is used to unmarshal the json response provided by call to
// Teamwork API endpoint /tasks/{id}/time/total.json.
type TaskTimeTotalsJSON struct {
	Data []*TaskTimeTotalJSON `json:"projects"`
}

// TaskQueryParams defines valid query parameters for this resource.
type TaskQueryParams struct {
	AssignedUserID   string `url:"responsible-party-ids,omitempty"`
	FromDate         string `url:"startDate,omitempty"`
	ToDate           string `url:"endDate,omitempty"`
	IncludeCompleted bool   `url:"includeCompletedTasks,omitempty"`
	Include          string `url:"include,omitempty"`
	ProjectIDs       string `url:"projectIds,omitempty"`
	PageSize         string `url:"pageSize,omitempty"`
	CompletedBefore  string `url:"completedBefore"`
	CompletedAfter   string `url:"completedAfter"`
}

func (resMsg *TaskResponseHandlerV3) ParseResponse(httpMethod string, rawRes []byte) error {
	// b := string(rawRes)
	// fmt.Println(b)

	err := json.Unmarshal(rawRes, &resMsg)
	if err != nil {
		return err
	}
	// fmt.Println("REPOSNE")
	// fmt.Println(resMsg.Response)

	if resMsg.Status == "Error" {
		return fmt.Errorf("received ERROR response: %s", resMsg.Message)
	}

	switch httpMethod {
	case http.MethodPost:
		// ABC€

		if resMsg.Task.ID == 0 {
			return fmt.Errorf("no task id returned for Task Post request ")
		}
	}

	return nil
}

// FormatQueryParams formats query parameters for this resource.
func (qp TaskQueryParams) FormatQueryParams() (string, error) {

	if qp.FromDate != "" {
		_, err := time.Parse("20060102", qp.FromDate)
		if err != nil {
			return "", fmt.Errorf("invalid format for FromDate parameter.  Should be YYYYMMDD, but found %s", qp.FromDate)
		}
	}

	if qp.ToDate != "" {
		_, err := time.Parse("20060102", qp.ToDate)
		if err != nil {
			return "", fmt.Errorf("invalid format for ToDate parameter.  Should be YYYYMMDD, but found %s", qp.ToDate)
		}
	}

	params, err := query.Values(qp)
	if err != nil {
		return "", err
	}

	return params.Encode(), nil
}

// GetTaskByID retrieves a specific task based on ID.
func (conn *Connection) GetTaskByIDV3(ID string) (*TaskVersion3, error) {

	_, err := strconv.Atoi(ID)
	if err != nil {
		if ID == "" {
			return nil, fmt.Errorf("missing required parameter(s): ID")
		}
		return nil, fmt.Errorf("invalid value (%s) for ID", ID)
	}

	endpoint := "tasks/" + ID
	data, err := conn.GetRequest(endpoint, nil)

	if err != nil {
		return nil, err
	}

	t := new(TaskVersion3)

	err = json.Unmarshal(data, &t)

	if err != nil {
		return nil, err
	}

	if t.Task.Id == 0 {
		return nil, fmt.Errorf("failed to retrieve task with ID (%s)", ID)
	}
	
	return t, nil
}
 
// GetTaskByID retrieves a specific task based on ID.
func (conn *Connection) GetTaskByID(ID string) (*Task, error) {

	_, err := strconv.Atoi(ID)
	if err != nil {
		if ID == "" {
			return nil, fmt.Errorf("missing required parameter(s): ID")
		}
		return nil, fmt.Errorf("invalid value (%s) for ID", ID)
	}

	endpoint := "tasks/" + ID

	data, err := conn.GetRequest(endpoint, nil)

	if err != nil {
		return nil, err
	}

	t := new(TaskJSON)

	err = json.Unmarshal(data, &t)

	if err != nil {
		return nil, err
	}

	if t.Task == nil {
		return nil, fmt.Errorf("failed to retrieve task with ID (%s)", ID)
	}

	return t.Task, nil
}

// GetTasks returns an array of tasks based on one or more query parameters.
func (conn *Connection) GetTasks(queryParams TaskQueryParams) ([]*Task, error) {

	data, err := conn.GetRequest("tasks", queryParams)
	if err != nil {
		return nil, err
	}

	tasks := new(TasksJSON)

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks.Tasks, nil
}


func (conn *Connection) PatchTask(taskID string, putData TaskPatchV3JSON)(int,error){
	handler := new(TaskResponseHandlerV3)

	b, err := json.Marshal(putData)
	if err != nil {
		return 0, err
	}
	
	err = conn.PatchRequest("tasks/"+taskID, b, handler)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	return handler.Task.ID, nil
}

// Creates a Task given the task list Id
func (conn *Connection) PostTask(taskListID string, postData TaskV3JSON) (int, error) {

	handler := new(TaskResponseHandlerV3)
	//b :=byteData.Bytes()
	//fmt.Printf(handler.Message)
	b, err := json.Marshal(postData)
	if err != nil {
		return 0, err
	}

	err1 := conn.PostRequest("tasklists/"+taskListID+"/tasks", b, handler)
	if err1 != nil {
		return 0, err1
	}

	if err != nil {
		return 0, err
	}

	return handler.Task.ID, nil
}

//Creates a subtask given the parent's task ID
func (conn *Connection) PostSubTask(parentTaskID string, postData TaskV3JSON) (int, error) {
	
	handler := new(TaskResponseHandlerV3)
	b, err := json.Marshal(postData)
	if err != nil {
		return 0, err
	}

	err1 := conn.PostRequest("tasks/"+parentTaskID+"/subtasks", b, handler)
	if err1 != nil {
		return 0, err1
	}

	if err != nil {
		return 0, err
	}

	return handler.Task.ID, nil
}

// GetTaskHours returns actual and estimated hours, and percent error in
// estimated hours for the specified task.
func (conn *Connection) GetTaskHours(taskID string) (*TimeTotals, error) {

	endpoint := fmt.Sprintf("tasks/%s/time/total", taskID)

	data, err := conn.GetRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}

	timeTotalsJSON := new(TaskTimeTotalsJSON)

	err = json.Unmarshal(data, &timeTotalsJSON)
	if err != nil {
		return nil, err
	}

	if len(timeTotalsJSON.Data) != 1 {
		return nil, fmt.Errorf("expected TaskTimeTotals to be size 1 but got %d", len(timeTotalsJSON.Data))
	}

	estimatedHours, err := strconv.ParseFloat(timeTotalsJSON.Data[0].Tasklist.Task.TimeEstimates.EstimatedHours, 64)
	if err != nil {
		return nil, err
	}

	actualHours, err := strconv.ParseFloat(timeTotalsJSON.Data[0].Tasklist.Task.TimeTotals.ActualHours, 64)
	if err != nil {
		return nil, err
	}

	return &TimeTotals{
		EstimatedHours: estimatedHours,
		ActualHours:    actualHours,
		PercentError:   CalculateEstimateError(estimatedHours, actualHours),
	}, nil
}

// CalculateEstimateError determines the percent error of a time estimate.
func CalculateEstimateError(estimate float64, actual float64) float64 {

	accuracy := (estimate - actual) / estimate * 100

	return math.Round(accuracy*100) / 100
}

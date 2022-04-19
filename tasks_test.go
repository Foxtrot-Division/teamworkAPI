package teamworkapi

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"
	"fmt"
)

// type tw struct{
// 	twAPIKey string `json:"tw_api_key"`
// }

type taskTestData struct {
	UserIDs       string `json:"users"`
	From          string `json:"from"`
	To            string `json:"to"`
	Count         int    `json:"count"`
	ExampleTaskID string `json:"exampleTaskID"`
	Include       string `json:"include"`
}

type parentTaskJson struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type CreateTaskTestData struct {
	Description      string         `json:"description"`
	EstimatedMinutes int            `json:"estimatedMinutes"`
	Name             string         `json:"name"`
	ParentTaskId     parentTaskJson `json:"parentTask"`
	Private          bool           `json:"private"`
}

type taskTestDataJSON struct {
	Data []taskTestData `json:"data"`
}

func initTaskTestConnectionV3(t *testing.T) *Connection {
	
	f, err := os.Open("./testdata/tw_api_conf.json")
	defer f.Close()

	if err != nil {
		t.Errorf(err.Error())
	}

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	data := new(TWAPIConf)

	err = json.Unmarshal(raw, &data)
	if err != nil {
		t.Errorf(err.Error())
	}

	//	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	conn, err := NewConnection(data.APIKey, data.SiteName, "", data.APIVersion)
	if err != nil {
		t.Fatalf(err.Error())
	}
	
	return conn
}

func loadTaskTestDataV3(t *testing.T) TaskV3JSON {

	f, err := os.Open("./testdata/taskTestData.json")
	defer f.Close()

	if err != nil {
		t.Errorf(err.Error())
	}

	data := new(TaskV3JSON)

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		t.Errorf(err.Error())
	}

	return *data
}

func loadSubTaskTestDataV3(t *testing.T) TaskV3JSON {

	f, err := os.Open("./testdata/subTaskTestData.json")
	defer f.Close()

	if err != nil {
		t.Errorf(err.Error())
	}

	data := new(TaskV3JSON)

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		t.Errorf(err.Error())
	}

	return *data
}

func initTaskTestConnection(t *testing.T) *Connection {
	//	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
	conn, err := NewConnection("twp_UtkI6MKeqjAgnW9hUtM8col7WTf7", "foxtrotdivision", "", "v3")
	if err != nil {
		t.Fatalf(err.Error())
	}

	return conn
}

func loadTaskTestData(t *testing.T) []taskTestData {

	f, err := os.Open("./testdata/taskTestData.json")
	defer f.Close()

	if err != nil {
		t.Errorf(err.Error())
	}

	data := new(taskTestDataJSON)

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

func TestFormatQueryParams(t *testing.T) {
	// test error responses
	var tests = []struct {
		qp   TaskQueryParams
		want url.Values
	}{
		{TaskQueryParams{AssignedUserID: "123456,102040"}, url.Values{"responsible-party-ids": {"123456,102040"}}},
		{TaskQueryParams{FromDate: "20201201", ToDate: "20201230"}, url.Values{"startDate": {"20201201"}, "endDate": {"20201230"}}},
		{TaskQueryParams{IncludeCompleted: true}, url.Values{"includeCompletedTasks": {"true"}}},
		{TaskQueryParams{}, url.Values{}},
	}

	for _, v := range tests {

		s, err := v.qp.FormatQueryParams()
		if err != nil {
			t.Errorf(err.Error())
		}

		result, err := url.ParseQuery(s)
		if err != nil {
			t.Errorf(err.Error())
		}

		if !reflect.DeepEqual(result, v.want) {
			t.Errorf("expected query keys/values (%s) but got (%s)", v.want, result)
		}
	}

	var tests2 = []struct {
		qp   TaskQueryParams
		want string
	}{
		{TaskQueryParams{FromDate: "10/2/2020", ToDate: "10/5/2020"}, "invalid format for FromDate parameter.  Should be YYYYMMDD, but found 10/2/2020"},
		{TaskQueryParams{FromDate: "20201002", ToDate: "10/5/2020"}, "invalid format for ToDate parameter.  Should be YYYYMMDD, but found 10/5/2020"},
	}

	for _, v := range tests2 {

		_, err := v.qp.FormatQueryParams()
		if err == nil {
			t.Errorf("expected error")
		}

		if err.Error() != v.want {
			t.Errorf("expected error (%s) but got (%s)", v.want, err.Error())
		}
	}
}

func TestGetTaskByIDV3(t *testing.T){
	taskId := "24064059"
	conn := initTaskTestConnectionV3(t)

	//task, err := conn.GetTaskByID(v.ExampleTaskID)
	task, err := conn.GetTaskByIDV3(taskId)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(task.Task.Attachments[0])
}

func TestGetSubtaskV3(t *testing.T){
	taskId := "24059040"
	conn := initTaskTestConnectionV3(t)

	//task, err := conn.GetTaskByID(v.ExampleTaskID)
	task, err := conn.GetSubtaskV3(taskId)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(task)
}

func TestGetTaskByID(t *testing.T) {

	testData := loadTaskTestData(t)

	conn := initTaskTestConnection(t)

	// test valid cases
	for _, v := range testData {

		//task, err := conn.GetTaskByID(v.ExampleTaskID)
		task, err := conn.GetTaskByID("14288573")
		if err != nil {
			t.Errorf(err.Error())
		}

		if task == nil {
			t.Errorf("no data returned for ID (%s)", v.ExampleTaskID)
		} else {
			id, err := strconv.Atoi(v.ExampleTaskID)
			if err != nil {
				t.Errorf(err.Error())
			}

			if task.ID != id {
				t.Errorf("expected ID (%d) but got (%d)", id, task.ID)
			}

			if task.Title == "" {
				t.Errorf("task id (%d) has no title", task.ID)
			}
		}
	}

	// test error responses
	var tests = []struct {
		ID   string
		want string
	}{
		{"123456", "failed to retrieve task with ID (123456)"},
		{"bad-content", "invalid value (bad-content) for ID"},
		{"", "missing required parameter(s): ID"},
	}

	for _, v := range tests {
		_, err := conn.GetTaskByID(v.ID)
		if err != nil {
			if err.Error() != v.want {
				t.Errorf("expected error (%s) but got (%s)", v.want, err.Error())
			}
		} else {
			t.Errorf("Expected error for taskID (%s)", v.ID)
		}
	}
}

func TestGetTasks(t *testing.T) {
	testData := loadTaskTestData(t)
	conn := initTaskTestConnection(t)

	for _, v := range testData {

		q := TaskQueryParams{
			AssignedUserID: v.UserIDs,
			FromDate:       v.From,
			ToDate:         v.To,
			Include:        v.Include,
		}

		tasks, err := conn.GetTasks(q)
		if err != nil {
			t.Errorf(err.Error())
		}

		if len(tasks) < 1 {
			t.Errorf("no tasks returned %s", conn.RequestURL)
		}
	}
}

func TestAttachFile(t *testing.T){
	conn := initTaskTestConnectionV3(t)

	f, err := os.Open("testdata/tw_api_conf.json")
	defer f.Close()

	if err != nil {
		t.Error(err)
	}

	data := new(TWAPIConf)

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		t.Error(err)
	}

	fc, err := NewFileConnection(data.SiteName, "Test_PDF_New.pdf", "testdata/Test_PDF_New.pdf", data.APIKey)
	if err != nil{
		t.Error(err)
	}

	fileID, err := fc.PutFile()
	if err != nil{
		t.Error(err)
	}

	jsonString := fmt.Sprintf(`{"attachments": {"pendingFiles": [{"categoryId": 0, "reference": "%v"}]}}`,fileID)
	//jsonString := fmt.Sprintf(`{"attachments":{"pendingFiles":[{"categoryId":0, "reference":"%v"}]}}`,"tf_341a8e1f-c840-440e-8891-ff3c70957643")

	//data := new(TaskPatchV3JSON)
	var taskPatchV3JSON TaskPatchV3JSON	

	err = json.Unmarshal([]byte(jsonString), &taskPatchV3JSON)
	if err != nil {
		t.Error(err)
	}

	ret, err := conn.PatchTask("24040057", taskPatchV3JSON)
	if err != nil{
		t.Error(err)
	}

	if ret == 0{
		t.Errorf("unknown error attaching file to task")
	}
}

func TestCreateTask(t *testing.T) {
	conn := initTaskTestConnectionV3(t)
	testData := loadTaskTestDataV3(t)

	fmt.Println(testData)

	_, err := conn.PostTask("1781185", testData)
	if err != nil {
		t.Errorf(err.Error())
	}
	//fmt.Println(ret)

}

func TestCreateSubTask(t *testing.T) {
	conn := initTaskTestConnectionV3(t)
//	testData := loadSubTaskTestDataV3(t)

	//taskJSON := fmt.Sprintf(`{"task": {"name": "Verify Time Logged/%v %v, %v","description": "%v","parentTaskId": %v,"assignees": {"userIds": [%v]}}}`,"2022-04-01", "Darth", "Vader", "TEST", 24040057, 179618)

	testJSON := `{"task": {"name": "Verify Time Logged/2022-04-01 Matt, Shilinski","description": "TEST","parentTaskId": 24050618,"assignees": {"userIds": [179618]}}}`

	var taskJSONData TaskV3JSON

	err := json.Unmarshal([]byte(testJSON), &taskJSONData)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(taskJSONData)

	ret, err := conn.PostSubTask("1781185", taskJSONData)
	if err != nil {
		t.Errorf(err.Error())
	}

	if ret == 0{
		t.Error("unknown error")
	}
	//	fmt.Println(postSubTask)

}

func TestGetTaskTotalHours(t *testing.T) {

	conn := initTaskTestConnection(t)

	var tests = []struct {
		taskID           string
		wantActual       float64
		wantEstimated    float64
		WantPercentError float64
	}{
		{"21603507", 222.30, 220.00, -1.05},
		{"21585386", 3.25, 25, 87},
	}

	for _, v := range tests {
		totals, err := conn.GetTaskHours(v.taskID)
		if err != nil {
			t.Errorf(err.Error())
		}

		if totals.ActualHours != v.wantActual {
			t.Errorf("expected actual hours to be %f but got %f", v.wantActual, totals.ActualHours)
		}

		if totals.EstimatedHours != v.wantEstimated {
			t.Errorf("expected actual hours to be %f but got %f", v.wantEstimated, totals.EstimatedHours)
		}

		if totals.PercentError != v.WantPercentError {
			t.Errorf("expected accuracy to be %f but got %f", v.WantPercentError, totals.PercentError)
		}
	}
}

func TestCalculateEstimateError(t *testing.T) {

	var tests = []struct {
		actual   float64
		estimate float64
		want     float64
	}{
		{222.30, 220.00, -1.05},
		{30, 5, -500},
		{5, 10, 50},
		{2, 2, 0},
	}

	for _, v := range tests {

		accuracy := CalculateEstimateError(v.estimate, v.actual)

		if accuracy != v.want {
			t.Errorf("expected percent error to be %f but got %f", v.want, accuracy)
		}
	}
}

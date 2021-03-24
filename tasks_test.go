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

type taskTestData struct {
	UserIDs 		string 	`json:"users"`
	From			string 	`json:"from"`
	To 				string 	`json:"to"`
	Count   		int		`json:"count"`
	ExampleTaskID 	string 	`json:"exampleTaskID"`
	Include 		string  `json:"include"`
}

type taskTestDataJSON struct {
	Data []taskTestData `json:"data"`
}

func initTaskTestConnection(t *testing.T) *Connection {
	conn, err := NewConnectionFromJSON("./testdata/apiConfigTestData1.json")
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
		qp 		TaskQueryParams
		want	url.Values
	}{
		{TaskQueryParams{AssignedUserID: "123456,102040"}, url.Values{"responsible-party-ids":{"123456,102040"}}},
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
		qp 		TaskQueryParams
		want	string
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
		ID 		string
		want	string
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

		q := TaskQueryParams {
			AssignedUserID: v.UserIDs,
			FromDate: v.From,
			ToDate: v.To,
			Include: v.Include,
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
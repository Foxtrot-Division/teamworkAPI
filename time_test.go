package teamworkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"
)

type TimeTestData struct {
	People      []string            `json:"people"`
	TaskID      string              `json:"taskID"`
	TimePeriods []map[string]string `json:"timePeriods"`
}

func initTimeTestConnection(t *testing.T) *Connection {
	conn, err := NewConnection("water589meat", "foxtrotdivision", "", "v3")
	if err != nil {
		t.Fatalf(err.Error())
	}
	//fmt.Print("Connection")
	return conn
}

func initTimeTestConnectionV3(t *testing.T) *Connection {
	
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

func initTimeTestData(t *testing.T) *TimeTestData {
	testData := new(TimeTestData)

	f, err := os.Open("./testdata/timeTestData.json")
	defer f.Close()

	if err != nil {
		t.Fatalf(err.Error())
	}

	raw, _ := ioutil.ReadAll(f)
	//fmt.Println(raw)
	err = json.Unmarshal(raw, &testData)
	if err != nil {
		t.Fatalf(err.Error())
	}

	return testData
}

func initTimeTestDataV3(t *testing.T) TimeQueryParamsV3 {

	f, err := os.Open("./testdata/timeTestData.json")
	defer f.Close()

	if err != nil {
		t.Errorf(err.Error())
	}

	data := new(TimeQueryParamsV3)

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		t.Errorf(err.Error())
	}

	//fmt.Println(data.ProjectID)

	return *data
}

// func initTimeTestDataV3(t *testing.T) *TimeQueryParamsV3 {
// 	testData := new(TimeQueryParamsV3)

// 	f, err := os.Open("./testdata/timeTestData.json")
// 	defer f.Close()

// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	raw, _ := ioutil.ReadAll(f)

// 	err = json.Unmarshal(raw, &testData)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	return testData
// }

func TestGetTimeEntries(t *testing.T) {

	conn := initTimeTestConnection(t)

	var tests = []struct {
		queryParams *TimeQueryParams
		want        int
	}{
		{&TimeQueryParams{FromDate: "20210101", ToDate: "20210228", PageSize: "500"}, 218},
		{&TimeQueryParams{FromDate: "20210304", ToDate: "20210304"}, 4},
	}

	for _, v := range tests {
		entries, err := conn.GetTimeEntries(v.queryParams)
		if err != nil {
			t.Errorf(err.Error())
		}
		if len(entries) != v.want {
			t.Errorf("expected %d time entries but got %d", v.want, len(entries))
		}

		for i := 1; i < len(entries); i++ {

			currentDate, err := time.Parse(TeamworkDateFormatLong, entries[i].Date)
			if err != nil {
				t.Errorf(err.Error())
			}
			priorDate, err := time.Parse(TeamworkDateFormatLong, entries[i-1].Date)
			if err != nil {
				t.Errorf(err.Error())
			}

			if priorDate.Equal(currentDate) {
				continue
			}
			if !priorDate.Before(currentDate) {
				t.Errorf("dates out of order - %s is not before %s", priorDate, currentDate)
			}
		}
	}
}

func TestGetTimeEntriesV3(t *testing.T) {

	conn := initTimeTestConnectionV3(t)

	timeQueryParamsV3 := TimeQueryParamsV3{
		AssignedToUserIds: []string{"266242"},
		EndDate: "2022-04-01",
		StartDate: "2022-03-18",
		PageSize: "500",
	}

	ret, err := conn.GetTimeEntriesV3(&timeQueryParamsV3)
	if err != nil {
		t.Errorf(err.Error())
	}

	//UserID, MIN, TASKID
	for _, u := range ret {
		fmt.Printf(" %v \n", *u)
	}
	// fmt.Print(len(ret))
	// fmt.Println(&ret)
}

func TestGetTimeEntriesByTask(t *testing.T) {

	conn := initTimeTestConnection(t)

	var tests = []struct {
		ID   string
		want int
	}{
		{"21603507", 14},
	}

	for _, v := range tests {
		entries, err := conn.GetTimeEntriesByTask(v.ID)
		if err != nil {
			t.Errorf(err.Error())

		}
		if len(entries) != v.want {
			t.Errorf("expected %d time entries for task %s but got %d", v.want, v.ID, len(entries))
		}
	}
}

func TestGetTimeEntriesByPerson(t *testing.T) {

	conn := initTimeTestConnection(t)

	testData := initTimeTestData(t)

	testDateLayout := "20060102"

	for _, p := range testData.People {
		for _, tp := range testData.TimePeriods {

			fromDate, err := time.Parse(testDateLayout, tp["fromdate"])
			if err != nil {
				t.Errorf(err.Error())
			}

			toDate, err := time.Parse(testDateLayout, tp["todate"])
			if err != nil {
				t.Errorf(err.Error())
			}

			entries, err := conn.GetTimeEntriesByPerson(p, tp["fromdate"], tp["todate"])

			if err != nil {
				t.Errorf(err.Error())
			}

			if len(entries) < 1 {
				t.Errorf("No time entries for person %s, from %s to %s.", p, tp["fromdate"], tp["todate"])
			}

			for _, entry := range entries {

				if entry.PersonID != p {
					t.Errorf("Found user ID (%s) but expected only (%s)", entry.PersonID, p)
				}

				entryTime, err := time.Parse(time.RFC3339, entry.Date)
				if err != nil {
					t.Errorf(err.Error())
				}
				d := time.Date(entryTime.Year(), entryTime.Month(), entryTime.Day(), 0, 0, 0, 0, time.UTC)

				if d.Before(fromDate) || d.After(toDate) {
					t.Errorf("Entry (%s) is not within specified time range (%s - %s)!", d, fromDate, toDate)
				}
			}
		}
	}

	var tests = []struct {
		ID   string
		from string
		to   string
		want string
	}{
		{"", "20201012", "20201013", "missing required parameter(s): personID"},
		{"abc", "", "", "strconv.Atoi: parsing \"abc\": invalid syntax"},
		{"12345", "10-12-2020", "20201013", "invalid format for FromDate parameter.  Should be YYYYMMDD, but found 10-12-2020"},
	}

	for _, v := range tests {
		_, err := conn.GetTimeEntriesByPerson(v.ID, v.from, v.to)
		if err != nil {
			if err.Error() != v.want {
				t.Errorf("expected error (%s) but got (%s)", v.want, err.Error())
			}
		} else {
			t.Errorf("Expected error for userID (%s)", v.ID)
		}
	}
}

func TestPostTimeEntry(t *testing.T) {

	conn := initTimeTestConnection(t)

	testData := initTimeTestData(t)

	type testCase struct {
		entry *TimeEntry
		error bool
		want  string
	}

	tests := make([]testCase, len(testData.People))

	for i, v := range testData.People {
		tests[i].entry = &TimeEntry{
			PersonID:    v,
			Description: fmt.Sprintf("test entry %d", i),
			Hours:       strconv.Itoa(5 + i),
			Minutes:     "0",
			Date:        time.Now().Format("20060102"),
			IsBillable:  "false",
			TaskID:      testData.TaskID,
		}
		tests[i].error = false
		tests[i].want = ""
	}

	tests = append(tests,
		testCase{entry: &TimeEntry{PersonID: "", Hours: "0", Minutes: "10", Date: "20201201", TaskID: ""}, error: true, want: "time entry is missing required field(s): PersonID, TaskID"},
		testCase{entry: &TimeEntry{PersonID: testData.People[0], Hours: "0", Minutes: "10", Date: "20201201", TaskID: "123456"}, error: true, want: "received ERROR response: Not Found"})

	for _, v := range tests {

		res, err := conn.PostTimeEntry(v.entry)

		if err != nil {
			if !v.error {
				t.Errorf(err.Error())
			} else {
				if err.Error() != v.want {
					t.Errorf("expected error (%s) but got (%s)", v.want, err.Error())
				}
			}
		} else {
			if v.error {
				t.Errorf("expected error")
			} else {
				if v.entry.ID != res {
					t.Errorf("ID (%s) not set to expected value (%s)", v.entry.ID, res)
				}
			}
		}

	}
}

func TestDeleteTimeEntry(t *testing.T) {

	conn := initTimeTestConnection(t)

	testData := initTimeTestData(t)

	testEntry := &TimeEntry{
		PersonID:    testData.People[0],
		Description: fmt.Sprintf("test entry - DELETE"),
		Hours:       "5",
		Minutes:     "0",
		Date:        time.Now().Format("20060102"),
		IsBillable:  "false",
		TaskID:      testData.TaskID,
	}

	id, err := conn.PostTimeEntry(testEntry)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Printf("test entry %s\n", id)
	err = conn.DeleteTimeEntry(id)
	if err != nil {
		t.Errorf(err.Error())
	}

	tests := []struct {
		ID    string
		error bool
		want  string
	}{
		{id, false, ""},
		{"", true, "missing required parameter: ID"},
		{"12345", true, "received ERROR response: Forbidden"},
		{"abc!def", true, "received ERROR response: Bad Request"},
	}

	for _, v := range tests {

		err := conn.DeleteTimeEntry(v.ID)

		if err != nil {
			if !v.error {
				t.Errorf(err.Error())
			} else {
				if err.Error() != v.want {
					t.Errorf("expected error (%s) but got (%s)", v.want, err.Error())
				}
			}
		} else {
			if v.error {
				t.Errorf("expected error for entry %s but got none", v.ID)
			}
		}
	}
}

func TestTotalAndAvgHours(t *testing.T) {

	var tests = []struct {
		entries   []*TimeEntry
		wantTotal float64
		wantAvg   float64
	}{
		{[]*TimeEntry{{Hours: "10", Minutes: "0"}, {Hours: "5", Minutes: "0"}}, 15.00, 7.50},
		{[]*TimeEntry{{Hours: "1", Minutes: "30"}, {Hours: "0", Minutes: "45"}}, 2.25, 1.13},
		{[]*TimeEntry{{Hours: "1", Minutes: "59"}, {Hours: "4", Minutes: "1"}}, 6, 3},
	}

	for _, v := range tests {

		r, err := TotalAndAvgHours(v.entries)
		if err != nil {
			t.Errorf(err.Error())
		}

		if r["total"] != v.wantTotal {
			t.Errorf("expected total hours to be %f but got %f", v.wantTotal, r["total"])
		}

		if r["avg"] != v.wantAvg {
			t.Errorf("expected avg hours to be %f but got %f", v.wantAvg, r["avg"])
		}
	}
}

func TestDurationInDays(t *testing.T) {

	var tests = []struct {
		from string
		days int
	}{
		{"20210101", 380},
	}

	for _, v := range tests {

		start, err := time.Parse(TeamworkDateFormatShort, v.from)
		if err != nil {
			t.Errorf(err.Error())
		}

		for i := 0; i < v.days; i++ {

			d := start.AddDate(0, 0, i)

			days, err := DurationInDays(start.Format(TeamworkDateFormatShort), d.Format(TeamworkDateFormatShort))
			if err != nil {
				t.Errorf(err.Error())
			}

			if days != i {
				t.Errorf("expected diff of %d day(s) between %s and %s but got %d", i,
					start.Format(TeamworkDateFormatShort), d.Format(TeamworkDateFormatShort), days)
			}
		}
	}

}

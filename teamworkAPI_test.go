package teamworkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const testDataDir string = "./testdata/"
const apiConfigData string = "apiConfig.json"

type APIConfig struct {
	APIKey         string `json:"apiKey"`
	SiteName       string `json:"siteName"`
	DataPreference string `json:"dataPreference"`
}

var apiConfig APIConfig

func TestMain(m *testing.T) {
	f, err := os.Open(testDataDir + apiConfigData)
	if err != nil {
		m.Fatalf("Failed to read test config file: %s", testDataDir+apiConfigData)
	}

	byteValue, _ := ioutil.ReadAll(f)

	json.Unmarshal(byteValue, &apiConfig)

	f.Close()
}

func TestNewTeamworkAPI(t *testing.T) {
	key := "someKey"
	site := "someSite"
	dataPref := "json"

	var api *teamworkAPI
	var err error

	api, err = NewTeamworkAPI(key, site, dataPref)

	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}

	if api.APIKey != key {
		t.Errorf("apiKey (%s) should be (%s)", api.APIKey, key)
	}

	if api.SiteName != site {
		t.Errorf("siteName (%s) should be (%s)", api.SiteName, site)
	}

	if api.DataPreference != dataPref {
		t.Errorf("dataPreference (%s) should be (%s)", api.DataPreference, dataPref)
	}

	api, err = NewTeamworkAPI(key, site, "")

	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}

	if api.DataPreference != "json" {
		t.Errorf("dataPreference (%s) should have defaulted to json", api.DataPreference)
	}

	api, err = NewTeamworkAPI("", "", "")

	if err == nil {
		t.Errorf("NewTeamworkAPI call was allowed with empty string parameter values.")
	}
}

func TestNewTeamworkAPIFromJSON(t *testing.T) {
	api, err := NewTeamworkAPIFromJSON(testDataDir + apiConfigData)

	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}

	if api.APIKey != apiConfig.APIKey {
		t.Errorf("apiKey (%s) should be (%s)", api.APIKey, apiConfig.APIKey)
	}

	if api.SiteName != apiConfig.SiteName {
		t.Errorf("siteName (%s) should be (%s)", api.SiteName, apiConfig.SiteName)
	}

	if api.DataPreference != apiConfig.DataPreference {
		t.Errorf("dataPreference (%s) should be (%s)", api.DataPreference, apiConfig.DataPreference)
	}

	if api.URL != "https://"+apiConfig.SiteName+".teamwork.com/" {
		t.Errorf("URL (%s) not formed correctly", api.URL)
	}
}

func TestGetRequest(t *testing.T) {
	api, _ := NewTeamworkAPIFromJSON(testDataDir + apiConfigData)

	_, err := api.GetRequest("projects." + api.DataPreference)

	if err != nil {
		t.Errorf(fmt.Sprintf("%s", err))
	}
}

package teamworkapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type teamworkAPI struct {
	APIKey         string `json:"apiKey"`
	SiteName       string `json:"siteName"`
	DataPreference string `json:"dataPreference"`
	URL            string
}

// NewTeamworkAPI initializes a new instance used to generate Teamwork API calls.
func NewTeamworkAPI(apiKey string, siteName string, dataPreference string) (*teamworkAPI, error) {

	var e string = ""

	if apiKey == "" {
		e += "apiKey\n"
	}
	if siteName == "" {
		e += "siteName"
	}

	if dataPreference == "" {
		dataPreference = "json"
	}

	if e != "" {
		return nil, errors.New("Missing required parameter(s):\n" + e)
	}

	t := new(teamworkAPI)
	t.APIKey = apiKey
	t.SiteName = siteName
	t.URL = "https://" + siteName + ".teamwork.com/"
	t.DataPreference = dataPreference

	return t, nil
}

// NewTeamworkAPIFromJSON initializes a new instance based on json file.
func NewTeamworkAPIFromJSON(pathToJSONFile string) (*teamworkAPI, error) {
	f, err := os.Open(pathToJSONFile)
	defer f.Close()
	if err != nil {
		return nil, errors.New("Failed to open JSON file at " + pathToJSONFile)
	}

	byteValue, _ := ioutil.ReadAll(f)

	t := new(teamworkAPI)

	json.Unmarshal(byteValue, &t)

	t.URL = "https://" + t.SiteName + ".teamwork.com/"

	return t, nil
}

// GetRequest performs a HTTP GET on the desired endpoint.
func (t teamworkAPI) GetRequest(endpoint string) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", t.URL+endpoint, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(t.APIKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func basicAuth(apiKey string) string {
	return base64.StdEncoding.EncodeToString([]byte(apiKey))
}

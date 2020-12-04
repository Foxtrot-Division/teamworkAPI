package teamworkapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type connection struct {
	APIKey         string `json:"apiKey"`
	SiteName       string `json:"siteName"`
	DataPreference string `json:"dataPreference"`
	URL            string
}

// NewConnection initializes a new instance used to generate Teamwork API calls.
func NewConnection(apiKey string, siteName string, dataPreference string) (*connection, error) {

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

	conn := new(connection)
	conn.APIKey = apiKey
	conn.SiteName = siteName
	conn.URL = "https://" + siteName + ".teamwork.com/"
	conn.DataPreference = dataPreference

	return conn, nil
}

// NewConnectionFromJSON initializes a new instance based on json file.
func NewConnectionFromJSON(pathToJSONFile string) (*connection, error) {
	f, err := os.Open(pathToJSONFile)
	defer f.Close()
	if err != nil {
		return nil, errors.New("Failed to open JSON file at " + pathToJSONFile)
	}

	byteValue, _ := ioutil.ReadAll(f)

	conn := new(connection)

	json.Unmarshal(byteValue, &conn)

	conn.URL = "https://" + conn.SiteName + ".teamwork.com/"

	return conn, nil
}

// GetRequest performs a HTTP GET on the desired endpoint.
func (conn connection) GetRequest(endpoint string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", conn.URL + endpoint + "." + conn.DataPreference, nil)

	req.Header.Add("Authorization", "Basic " + basicAuth(conn.APIKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func basicAuth(apiKey string) string {
	return base64.StdEncoding.EncodeToString([]byte(apiKey))
}
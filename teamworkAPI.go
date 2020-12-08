package teamworkapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

// GetRequest performs a HTTP GET on the desired endpoint, with the specific query parameters.
func (conn connection) GetRequest(endpoint string, params map[string] interface{}) ([]byte, error) {
	client := &http.Client{}

	queryParams := ""
	var err error

	if params != nil {
		query, err := FormatQueryString(params)

		if err != nil {
			return nil, err
		}

		queryParams += "?" + query.Encode()
	}

	req, err := http.NewRequest("GET", conn.URL + endpoint + "." + conn.DataPreference + queryParams, nil)

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

// FormatQueryString generates a http query string from the supplied map containing request parameters.
func FormatQueryString(params map[string] interface {}) (url.Values, error) {

	queryString := url.Values{}
	if params != nil {
		for key, value := range params {
			switch value.(type) {
			case string:
				queryString.Add(key, fmt.Sprintf("%s",value))
			
			case int:
				queryString.Add(key, fmt.Sprintf("%d", value))
			
			case bool:
				queryString.Add(key, fmt.Sprintf("%t", value))
			
			default:
				log.Printf("Unsupported type (%T) for %s.\n", value, key)
			}
		}
	}

	return queryString, nil
}

func basicAuth(apiKey string) string {
	return base64.StdEncoding.EncodeToString([]byte(apiKey))
}


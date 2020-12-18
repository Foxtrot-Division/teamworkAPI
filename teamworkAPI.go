// Package teamworkapi provides utilities to interface with the Teamwork API.
package teamworkapi

import (
	"bytes"
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

// Connection stores info needed to establish Teamwork API Connection
type Connection struct {
	APIKey         string `json:"apiKey"`
	SiteName       string `json:"siteName"`
	DataPreference string `json:"dataPreference"`
	URL            string
}

// NewConnection initializes a new instance used to generate Teamwork API calls.
// If dataPreference is empty string (""), it will default to json.
func NewConnection(apiKey string, siteName string, dataPreference string) (*Connection, error) {

	errBuff := ""

	if apiKey == "" {
		errBuff += "apiKey"
	}
	if siteName == "" {
		if errBuff != "" {
			errBuff += ", "
		}
		errBuff += "siteName"
	}

	if errBuff != "" {
		return nil, fmt.Errorf("missing required parameter(s): %s", errBuff)
	}

	if dataPreference == "" {
		dataPreference = "json"
	}

	conn := new(Connection)
	conn.APIKey = apiKey
	conn.SiteName = siteName
	conn.URL = "https://" + siteName + ".teamwork.com/"
	conn.DataPreference = dataPreference

	return conn, nil
}

// NewConnectionFromJSON initializes a new instance based on json file.
func NewConnectionFromJSON(pathToJSONFile string) (*Connection, error) {

	f, err := os.Open(pathToJSONFile)

	defer f.Close()
	
	if err != nil {
		return nil, errors.New("Failed to open JSON file at " + pathToJSONFile)
	}

	byteValue, _ := ioutil.ReadAll(f)

	conn := new(Connection)

	json.Unmarshal(byteValue, &conn)

	errBuff := ""

	if conn.APIKey == "" {
		errBuff += "apiKey"
	}

	if conn.SiteName == "" {
		if errBuff != "" {
			errBuff += ", "
		}
		errBuff += "siteName"
	}

	if errBuff != "" {
		return nil, fmt.Errorf("missing required parameter(s): %s", errBuff)
	}

	if conn.DataPreference == "" {
		conn.DataPreference = "json"
	}

	conn.URL = "https://" + conn.SiteName + ".teamwork.com/"

	return conn, nil
}

// GetRequest performs a HTTP GET on the desired endpoint, with the specific query parameters.
func (conn Connection) GetRequest(endpoint string, params map[string]interface{}) ([]byte, error) {
	
	if endpoint == "" {
		return nil, fmt.Errorf("missing required parameter(s): endpoint")
	}

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

	req, err := http.NewRequest("GET", conn.URL+endpoint+"."+conn.DataPreference+queryParams, nil)

	req.Header.Add("Authorization", "Basic "+basicAuth(conn.APIKey))

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

// PostRequest submits a POST request to Teamwork API.  It is up to the caller
// to properly marshal json into the data parameter.
func (conn *Connection) PostRequest(endpoint string, data []byte) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest("POST", conn.URL+endpoint+".json", bytes.NewBuffer(data))
	req.Header.Add("Authorization", "Basic "+basicAuth(conn.APIKey))
	req.Header.Set("Content-Type", "application/json")

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

// FormatQueryString generates a http query string from the supplied map containing request parameters.
func FormatQueryString(params map[string]interface{}) (url.Values, error) {

	queryString := url.Values{}
	if params != nil {
		for key, value := range params {
			switch value.(type) {
			case string:
				queryString.Add(key, fmt.Sprintf("%s", value))

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

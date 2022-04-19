package teamworkapi

import(
	"testing"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	//"strconv"
)

func initFileVersionTestConnectionV3(t *testing.T) *Connection {

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

func TestPostNewFileVersion(t *testing.T){

	conn := initFileVersionTestConnectionV3(t)
	jsonString := fmt.Sprintf(`{"fileversion": {"pendingFileRef":"%v","categoryId": %v}}`, "tf_83caff54-b376-4635-ae5e-2d5e7a1af8f3",858993)

	var fileVersionBody FileVersionBody

	err := json.Unmarshal([]byte(jsonString), &fileVersionBody)
	if err != nil {
		t.Error(err)
	}

	ret, err := conn.PostNewFileVersion("6673830", fileVersionBody)
	if err != nil{
		t.Error(err)
	}

	fmt.Println(ret)
}
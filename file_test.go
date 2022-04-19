package teamworkapi

import(
	"testing"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	//"strconv"
)

func initFileTestConnectionV3(t *testing.T) *Connection {

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

func TestPatchFile(t *testing.T){

	conn := initTaskTestConnectionV3(t)
	jsonString := `{"file": {"categoryId": 858993}}`
	
	 var file FileVersion3

	err := json.Unmarshal([]byte(jsonString), &file)
	if err != nil {
		t.Error(err)
	}

	ret, err := conn.PatchFile("6673096", file)
	if err != nil{
		t.Error(err)
	}

	fmt.Println(ret)
}

func TestPutFile(t* testing.T){
	f, err := os.Open("testdata/tw_api_conf.json")
	defer f.Close()

	if err != nil {
		t.Errorf(err.Error())
	}

	data := new(TWAPIConf)

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		t.Errorf(err.Error())
	}

	fc, err := NewFileConnection(data.SiteName, "Test_PDF_New.pdf", "/Users/matthewshilinski/Documents/Test_PDF_New.pdf", data.APIKey)
	if err != nil{
		t.Error(err)
	}

	// preSignedRes, err := fc.GetPreSignedData()
	// if err != nil{
	// 	t.Error(err)
	// }

	fileRef, err := fc.PutFile()
	if err != nil{
		t.Error(err)
	}

	fmt.Println(fileRef)
}
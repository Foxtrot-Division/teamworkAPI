package teamworkapi

import(
	"testing"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	//"strconv"
)


// func TestGetPreSignedURL(t *testing.T){

// 	f, err := os.Open("testdata/tw_api_conf.json")
// 	defer f.Close()

// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}

// 	data := new(TWAPIConf)

// 	raw, err := ioutil.ReadAll(f)
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}

// 	err = json.Unmarshal(raw, &data)
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}

// 	fi, err := os.Stat("testdata/Test_PDF_New.pdf")
//     if err != nil {
//         t.Error(err)
// 	}

// 	size := strconv.FormatInt(fi.Size(), 10)

// 	fc, err := NewFileConnection(data.SiteName, "Test_PDF_New.pdf", size, "testdata/Test_PDF_New.pdf", data.APIKey)
// 	if err != nil{
// 		t.Error(err)
// 	}

// 	preSignedRes, err := fc.GetPreSignedData()
// 	if err != nil{
// 		t.Error(err)
// 	}

// 	fmt.Println(preSignedRes)
// }

func TestPutFile(t*testing.T){
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
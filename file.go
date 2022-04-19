package teamworkapi

import(
	"fmt"
	"strings"
	"strconv"
	"io/ioutil"
	"io"
	"net/http"
	"mime/multipart"
	"bytes"
	"os"
	"github.com/imroc/req/v3"
	"encoding/json"
)

//models a file from TW, add addional fields as needed
type FileVersion3 struct{
	File struct{
		Id         int `json:"id,omitempty"`
		CategoryId int `json:"categoryId,omitempty"`
	} `json:"file"`
}

type FileConnection struct{
	SiteName        string
	FileName        string 
	FullPathToFile  string
	APIKey          string
}

type PreSignedRes struct{
	Ref string `json:"ref"`
	URL string `json:"url"`
}

type FileResponse struct{
	Id         int `json:'id'`
	CategoryId int `json:'categoryId'`
}

type FileResponseHandlerV3 struct {
	Status  string         `json:"STATUS"`
	Message string         `json:"MESSAGE"`
	File    FileResponse   `json:"file`
}
// // TaskResponseHandlerV3 models a http response for a Task operation using version 3 of teamwork api.
// type FileResponseHandlerV3 struct {
// 	File FileVersion3 `json:"file`
// }

// type FileResponseV3 struct {
// 	ID int `json:"id"`
// }


func (resMsg *FileResponseHandlerV3) ParseResponse(httpMethod string, rawRes []byte) error {
	// b := string(rawRes)
	// fmt.Println(b)

	err := json.Unmarshal(rawRes, &resMsg)
	if err != nil {
		return err
	}
	// fmt.Println("REPOSNE")
	// fmt.Println(resMsg.Response)

	if resMsg.Status == "Error" {
		return fmt.Errorf("received ERROR response: %s", resMsg.Message)
	}

	switch httpMethod {
	case http.MethodPost:
		// ABC€

		if resMsg.File.Id == 0 {
			return fmt.Errorf("no task id returned for Task Post request ")
		}
	}

	return nil
}

//Used for initially uploading a file
func NewFileConnection(SiteName string, FileName string, FullPathToFile string, APIKey string)(*FileConnection, error){

	if len(strings.TrimSpace(SiteName)) == 0{
		return  nil, fmt.Errorf("Missing SiteName value")
	}

	if len(strings.TrimSpace(FileName)) == 0{
		return nil, fmt.Errorf("Missing FileName value")
	}

	if len(strings.TrimSpace(FullPathToFile)) == 0{
		return nil, fmt.Errorf("Missing FullPathToFile value")
	}

	if len(strings.TrimSpace(APIKey)) == 0{
		return nil, fmt.Errorf("Missing APIKey value")
	}


	fc := new(FileConnection)
	fc.SiteName = SiteName
	fc.FileName = FileName
	fc.FullPathToFile = FullPathToFile
	fc.APIKey = APIKey

	return fc, nil
} 


//Specific Get Request to return the unique ref ID for said file and unique URL to PUT file to
func getPreSignedData(SiteName string, FileName string, ContentLength string, APIKey string) (*PreSignedRes, error){

	url := "https://" + SiteName + ".teamwork.com/projects/api/v1/pendingfiles/presignedurl.json?fileName=" + FileName + "&fileSize=" + ContentLength
	var preSignedRes *PreSignedRes

	client := req.C().DevMode()
	resp, err := client.R().
		SetBasicAuth(APIKey, "p").
		SetResult(&preSignedRes).
		Get(url)

	if err != nil{
		return nil, err
	}

	if !resp.IsSuccess(){
		return nil, fmt.Errorf("presigned URL request unsuccessful")
	}

	return preSignedRes, nil
}

func (fc *FileConnection) PutFile() (string, error){

	filePath := fc.FullPathToFile
	file, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filePath)
	io.Copy(part, file)
	writer.Close()

	//file.stat does not produce correct content length needed. Need actual req.body content length. 
	//Because the TW API requires the fileSize (really its the req.body size) when generating the preSignedUrl, 
	//we create a dummy request in order to get the actual req.body size. We then have the correct size/content-length 
	//to generated the preSignedURL (content length is used to during signature creation) for "putting" file to.
	req, err := http.NewRequest("PUT", "https://www.google.com", body)
	if err != nil{
		return "", err
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil{
		return "", err
	}

	contentLength := strconv.Itoa(len(reqBody))
	preSignedData, err := getPreSignedData(fc.SiteName, fc.FileName, contentLength, fc.APIKey)


	//Because we have read off the request body above to figure out the content length, need to re-create the multipart data
	filePath = fc.FullPathToFile
	file, _ = os.Open(filePath)
	defer file.Close()

	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	part, _ = writer.CreateFormFile("file", filePath)
	io.Copy(part, file)
	writer.Close()

	r, err := http.NewRequest("PUT", preSignedData.URL, body)
	if err != nil{
		return "", err
	}

	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.Header.Add("X-Amz-Acl", "public-read")
	r.Header.Add("Content-Length", contentLength)
	r.Header.Add("Host", "tw-bucket.s3-accelerate.amazonaws.com")

	client := &http.Client {}
	
	rsp, _ := client.Do(r)
	fmt.Println(rsp)
    if rsp.StatusCode != 200 {
        return "", fmt.Errorf("Request failed with response code: %d", rsp.StatusCode)
	}

	return preSignedData.Ref, nil
}

func (conn *Connection) PatchFile(fileID string, patchData FileVersion3) (*FileResponseHandlerV3, error) {

	handler := new(FileResponseHandlerV3)

	b, err := json.Marshal(patchData)
	if err != nil {
		return nil, err
	}

	err = conn.PatchRequest("files/"+fileID, b, handler)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return handler, nil
}
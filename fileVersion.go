package teamworkapi

import(
//	"fmt"
	"encoding/json"
)

type FileVersionBody struct{
	FileVersion struct{
		CategoryId        int `json:"categoryId,omitempty"`
		PendingFileRef    string `json:"pendingFileRef,omitempty"`
	} `json:"fileversion"`
}

type FileVersionRes struct{
	FileVersion struct{
		FileVersionId   int    `json:"fileVersionId"`
		Status     		string `json:"STATUS"`
		Message    		string `json:"MESSAGE"`
	} `json:"fileversion"`
}

func (resMsg *FileVersionRes) ParseResponse(httpMethod string, rawRes []byte) error {

	err := json.Unmarshal(rawRes, &resMsg)
	if err != nil {
		return err
	}
	return nil
}

func (conn *Connection) PostNewFileVersion(existingFileID string, postData FileVersionBody) (*FileVersionRes, error) {

	handler := new(FileVersionRes)

	b, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}

	err = conn.PostRequest("files/"+existingFileID, b, handler)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return handler, nil
}
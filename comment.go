package teamworkapi

import (
	"encoding/json"
	"fmt"
)

//Notify string "all" - means notify all project users. Notify "true"is for only followers.

type CommentJSON struct {
	Body        string `json:"body"`
	ContentType string `json:"content-type"`
	Notify      string `json:"notify"`
}
type Comment struct {
	Comment CommentJSON `json:"comment"`
}

type CommentResponseHandler struct {
	Status  string `json:"STATUS"`
	Message string `json:"MESSAGE"`
}

// type ResponseHandler struct {
// 	Status  string `json:"STATUS"`
// 	Message string `json:"MESSAGE"`
// }
func (resMsg *CommentResponseHandler) ParseResponse(httpMethod string, rawRes []byte) error {

	err := json.Unmarshal(rawRes, &resMsg)
	if err != nil {
		return err
	}

	if resMsg.Status == "Error" {
		return fmt.Errorf("received ERROR response: %s", resMsg.Message)
	}

	// switch httpMethod {
	// case http.MethodPost:
	// 	// if resMsg.ID == "" {
	// 	// 	return fmt.Errorf("no ID returned for Comment entry POST")
	// 	// }
	// }

	return nil
}
func (conn *Connection) PostComment(ResourceId string, postData CommentJSON) (string, error) {
	//resource can be links, milestones, files, notebooks or tasks
	handler := new(CommentResponseHandler)
	//b :=byteData.Bytes()
	//fmt.Printf(handler.Message)

	commentEntryJSON := Comment{
		Comment: postData,
	}
	data, err := json.Marshal(commentEntryJSON)
	if err != nil {
		return "", err
	}
	err1 := conn.PostRequest("tasks/"+ResourceId+"/comments", data, handler)
	if err1 != nil {
		return "", err1
	}

	return handler.Status, nil
}

package teamworkapi

import (
	"fmt"
	"testing"
)

func initCommentTestConnection(t *testing.T) *Connection {
	conn, err := NewConnection("", "foxtrotdivision", "", "v1")
	if err != nil {
		t.Fatalf(err.Error())
	}
	//fmt.Print("Connection")
	return conn
}

func TestPostComment(t *testing.T) {

	conn := initCommentTestConnection(t)

	q := CommentJSON{
		Body:        "Test adding comment",
		ContentType: "TEXT",
		Notify:      "",
	}

	events, err := conn.PostComment("23948144", q)
	if err != nil {
		fmt.Println("err")

		t.Fatalf(err.Error())
	}

	fmt.Println(events)
}

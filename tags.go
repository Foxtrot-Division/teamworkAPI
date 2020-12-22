package teamworkapi

import "encoding/json"

// Tag models an individual tag in Teamwork.
type Tag struct {
	IDBuff	interface{}	`json:"id"`		// ID can show up as int or string in API response
	Name	string 		`json:"name"`
	Color	string		`json:"color"`
}

// TagJSON provides a wrapper around Tag to properly marshal json
// data when posting to API.
type TagJSON struct {
	Tag	*Tag `json:"tag"`
}

// TagsJSON models the parent JSON structure of an array of Tags and
// facilitates unmarshalling.
type TagsJSON struct {
	Tags []*Tag `json:"tags"`
}

// GetTags gets all tags.
func (conn Connection) GetTags() ([]*Tag, error) {

	data, err := conn.GetRequest("tags", nil)
	if err != nil {
		return nil, err
	}

	raw := new(TagsJSON)

	err = json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}

	return raw.Tags, nil
}
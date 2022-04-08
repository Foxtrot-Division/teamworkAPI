package teamworkapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
)

// CalendarEvent models a Teamwork Calendar event.
type CalendarEvent struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Start       string     `json:"start"`
	End         string     `json:"end"`
	AllDay      bool       `json:"all-day"`
	Type        *EventType `json:"type"`
	AttendeeIDs string     `json:"attending-user-ids"`
	Status      string     `json:"status"`
}

// CalendarEventsJSON models the parent JSON structure of an array of CalendarEvent and
// facilitates unmarshalling.
type CalendarEventsJSON struct {
	Events []*CalendarEvent `json:"events"`
}

// EventType models a Teamwork Calendar Event Type.
type EventType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// CalendarEventResponseHandler models a http response for a Calendar Event operation.
type CalendarEventResponseHandler struct {
	Status  string `json:"STATUS"`
	Message string `json:"MESSAGE"`
}

// CalendarEventQueryParams defines valid query parameters for this resource.
type CalendarEventQueryParams struct {
	UserID      string `url:"userId,omitempty"`
	From        string `url:"startdate,omitempty"`
	To          string `url:"endDate,omitempty"`
	EventTypeID string `url:"eventTypeId,omitempty"`
}

type CalendarEventQueryParamsV3 struct {
	CreatedByUserID int    `url:"createdByUserId,omitempty"`
	StartDate       string `url:"startDate,omitempty"`
	EndDate         string `url:"endDate,omitempty"`
}

// FormatQueryParams formats query parameters for this resource.
func (qp CalendarEventQueryParams) FormatQueryParams() (string, error) {

	if qp.From != "" {
		_, err := time.Parse("20060102", qp.From)
		if err != nil {
			return "", fmt.Errorf("invalid format for From parameter.  Should be YYYYMMDD, but found %s", qp.From)
		}
	} else {
		return "", fmt.Errorf("missing required parameter 'From'")
	}

	if qp.To != "" {
		_, err := time.Parse("20060102", qp.To)
		if err != nil {
			return "", fmt.Errorf("invalid format for To parameter.  Should be YYYYMMDD, but found %s", qp.To)
		}
	}

	params, err := query.Values(qp)
	if err != nil {
		return "", err
	}

	return params.Encode(), nil
}

// GetCalendarEvents returns an array of tasks based on one or more query parameters.
func (conn *Connection) GetCalendarEvents(queryParams CalendarEventQueryParams) ([]*CalendarEvent, error) {

	data, err := conn.GetRequest("calendarevents", queryParams)
	if err != nil {
		return nil, err
	}

	events := new(CalendarEventsJSON)

	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, err
	}

	return events.Events, nil
}

// func (conn *Connection) GetCalendarEventsV3(queryParams CalendarEventQueryParamsV3) ([]*CalendarEvent, error) {

// 	data, err := conn.GetRequestV3("calendar/events", queryParams)
// 	if err != nil {
// 		return nil, err
// 	}

// 	events := new(CalendarEventsJSON)

// 	err = json.Unmarshal(data, &events)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return events.Events, nil
// }

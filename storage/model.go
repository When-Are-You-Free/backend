package storage

import (
	"errors"
	"time"
)

const dateFormat = "2006-01-02"

type Date struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in "YYYY-MM-DD" format.
func (t Date) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(dateFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, dateFormat)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in "YYYY-MM-DD" format.
func (t *Date) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	innerTime, err := time.Parse(`"`+dateFormat+`"`, string(data))
	t.Time = innerTime
	return err
}

type InvitedPerson struct {
	UserToken string `json:"-"`

	Name           string `json:"name"`
	AvailableDates []Date `json:"available_dates,omitempty"`
	PreferredDate  *Date  `json:"preferred_date,omitempty"`
}

type Meetup struct {
	Description    string           `json:"description"`
	From           Date             `json:"from"`
	To             *Date            `json:"to,omitempty"`
	InvitedPeople  []*InvitedPerson `json:"invited_person"`
	SuggestedDates []Date           `json:"suggested_dates,omitempty"`
	FinalDate      *Date            `json:"final_date,omitempty"`
	Locked         bool             `json:"locked"`
}

// HasAccess checks whether any person has the given token, indicating the user
// should have access.
func (m *Meetup) HasAccess(userToken string) bool {
	for _, person := range m.InvitedPeople {
		if person.UserToken == userToken {
			return true
		}
	}

	return false
}

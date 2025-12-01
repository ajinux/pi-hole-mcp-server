package client

import (
	"encoding/json"
	"time"
)

// UnixTime is a custom type for unmarshaling Unix timestamps
type UnixTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for Unix timestamps
func (ut *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return err
	}
	ut.Time = time.Unix(timestamp, 0)
	return nil
}

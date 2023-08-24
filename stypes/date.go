package stypes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const DatetimeFormat = "2006-01-02 15:04:05"
const DatetimeFormatShort = "2006-01-02 15:04"
const DateFormat = "2006-01-02"
const DateFormatShort = "2006-01"
const TimeFormat = "15:04:05"
const TimeFormatShort = "15:04"
const TimeZoneFormat = "Z07:00"

func ParseTime(t string) (time.Time, error) {
	//time.RFC3339: "2006-01-02T15:04:05Z07:00"
	hasZ := strings.Contains(t, "Z")
	hasT := strings.Contains(t, "T")
	hasPlus := strings.Contains(t, "+")
	hasTimeZone := hasZ || hasPlus
	dt := t
	if hasTimeZone {
		i := strings.Index(t, "+")
		if i == -1 {
			i = strings.Index(t, "Z")
		}
		dt = t[:i]
	}
	colonCount := strings.Count(dt, ":")
	dashCount := strings.Count(dt, "-")

	if hasT && hasPlus && colonCount == 2 && dashCount == 2 {
		return time.Parse(time.RFC3339, t)
	}
	if hasT && hasZ && colonCount == 2 && dashCount == 2 {
		return time.Parse(time.RFC3339, t)
	}
	if hasT {
		t = strings.Replace(t, "T", " ", 1)
	}

	var format string
	if colonCount == 2 && dashCount == 2 {
		format = DatetimeFormat
	} else if colonCount == 1 && dashCount == 2 {
		format = DatetimeFormatShort
	} else if colonCount == 0 && dashCount == 2 {
		format = DateFormat
	} else if colonCount == 0 && dashCount == 1 {
		format = DateFormatShort
	} else if colonCount == 2 && dashCount == 0 {
		format = TimeFormat
	} else if colonCount == 1 && dashCount == 0 {
		format = TimeFormatShort
	}
	if hasTimeZone {
		format += TimeZoneFormat
	}
	return time.Parse(format, t)
}

type Time struct {
	sql.NullTime
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (t *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.Valid = false
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	d, err := ParseTime(s)
	if err != nil {
		return err
	}
	t.Valid = true
	t.Time = d
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// It returns an empty string if invalid, otherwise time.Time's MarshalText.
func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return t.Time.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It has backwards compatibility with v3 in that the string "null" is considered equivalent to an empty string
// and unmarshaling will succeed. This may be removed in a future version.
func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	// allowing "null" is for backwards compatibility with v3
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	d, err := ParseTime(string(text))
	if err != nil {
		return err
	}
	t.Valid = true
	t.Time = d
	return nil

}

type Date struct {
	sql.NullTime
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Date) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(DateFormat))), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (t *Date) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.Valid = false
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	d, err := ParseTime(s)
	if err != nil {
		return err
	}
	t.Valid = true
	t.Time = d
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// It returns an empty string if invalid, otherwise time.Time's MarshalText.
func (t Date) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return t.Time.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It has backwards compatibility with v3 in that the string "null" is considered equivalent to an empty string
// and unmarshaling will succeed. This may be removed in a future version.
func (t *Date) UnmarshalText(text []byte) error {
	str := string(text)
	// allowing "null" is for backwards compatibility with v3
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	d, err := ParseTime(string(text))
	if err != nil {
		return err
	}
	t.Valid = true
	t.Time = d
	return nil
}

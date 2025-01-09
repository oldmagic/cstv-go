package util

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// Regular expressions for parsing GOTV+ tokens
var (
	steamIDExp   = regexp.MustCompile(`^s(\d{17,})`)
	timeStampExp = regexp.MustCompile(`t(\d{10,})$`)
	tokenExp     = regexp.MustCompile(`^s\d{17,}t\d{10,}$`)
)

// ErrInvalidToken is returned when the token format is incorrect.
var ErrInvalidToken = errors.New("invalid GOTV+ token format")

// ParseToken extracts the SteamID and timestamp from a GOTV+ token.
// Example input: "s845489096165654t8799308478907"
// Returns: ("845489096165654", time.Time{}, nil) or an error.
func ParseToken(token string) (string, time.Time, error) {
	if !tokenExp.MatchString(token) {
		return "", time.Time{}, ErrInvalidToken
	}

	// Extract SteamID
	matches := steamIDExp.FindStringSubmatch(token)
	if len(matches) < 2 {
		return "", time.Time{}, ErrInvalidToken
	}
	steamID := matches[1]

	// Extract Timestamp
	matches = timeStampExp.FindStringSubmatch(token)
	if len(matches) < 2 {
		return "", time.Time{}, ErrInvalidToken
	}
	ts, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return "", time.Time{}, err
	}

	return steamID, time.Unix(ts, 0), nil
}

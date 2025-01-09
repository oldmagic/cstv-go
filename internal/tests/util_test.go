package tests

import (
	"testing"
	"time"

	"github.com/FlowingSPDG/gotv-plus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestParseToken_Valid(t *testing.T) {
	token := "s845489096165654t1672531199"
	steamID, timestamp, err := util.ParseToken(token)

	assert.NoError(t, err)
	assert.Equal(t, "845489096165654", steamID)
	assert.Equal(t, time.Unix(1672531199, 0), timestamp)
}

func TestParseToken_InvalidFormat(t *testing.T) {
	token := "invalid_token"
	_, _, err := util.ParseToken(token)
	assert.Error(t, err)
}

func TestParseToken_MalformedSteamID(t *testing.T) {
	token := "sABCt1672531199"
	_, _, err := util.ParseToken(token)
	assert.Error(t, err)
}

func TestParseToken_MalformedTimestamp(t *testing.T) {
	token := "s845489096165654tINVALID"
	_, _, err := util.ParseToken(token)
	assert.Error(t, err)
}

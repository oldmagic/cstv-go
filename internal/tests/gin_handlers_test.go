package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FlowingSPDG/gotv-plus-go/gotv"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock implementations for testing
type MockAuth struct{}
type MockStore struct{}
type MockBroadcaster struct{}

func (m *MockAuth) Auth(token, auth string) error {
	if token != "valid_token" || auth != "valid_auth" {
		return gotv.ErrInvalidAuth
	}
	return nil
}

func (m *MockStore) OnFull(token string, fragment int, tick int, at time.Time, r io.Reader) error {
	if token == "not_found" {
		return gotv.ErrMatchNotFound
	}
	return nil
}

func (m *MockBroadcaster) GetSync(token string, fragment int) (gotv.Sync, error) {
	if token == "not_found" {
		return gotv.Sync{}, gotv.ErrMatchNotFound
	}
	return gotv.Sync{Tick: 100, Fragment: 5}, nil
}

func (m *MockBroadcaster) GetSyncLatest(token string) (gotv.Sync, error) {
	return gotv.Sync{Tick: 200, Fragment: 10}, nil
}

func TestCheckAuthMiddleware_ValidAuth(t *testing.T) {
	router := gin.Default()
	auth := &MockAuth{}
	ginCSTV := gotv.NewGinCSTV(auth, nil, nil)
	router.Use(ginCSTV.CheckAuthMiddleware())

	req, _ := http.NewRequest("GET", "/valid_token", nil)
	req.Header.Set("X-Origin-Auth", "valid_auth")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCheckAuthMiddleware_InvalidAuth(t *testing.T) {
	router := gin.Default()
	auth := &MockAuth{}
	ginCSTV := gotv.NewGinCSTV(auth, nil, nil)
	router.Use(ginCSTV.CheckAuthMiddleware())

	req, _ := http.NewRequest("GET", "/valid_token", nil)
	req.Header.Set("X-Origin-Auth", "invalid_auth")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOnStartFragment_ValidRequest(t *testing.T) {
	router := gin.Default()
	store := &MockStore{}
	ginCSTV := gotv.NewGinCSTV(nil, store, nil)

	router.POST("/:token/:fragment_number/start", ginCSTV.OnStartFragment())

	req, _ := http.NewRequest("POST", "/valid_token/1/start?tps=64&tick=120", bytes.NewBuffer([]byte("test")))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOnStartFragment_InvalidQuery(t *testing.T) {
	router := gin.Default()
	store := &MockStore{}
	ginCSTV := gotv.NewGinCSTV(nil, store, nil)

	router.POST("/:token/:fragment_number/start", ginCSTV.OnStartFragment())

	req, _ := http.NewRequest("POST", "/valid_token/1/start?tps=not_a_number", bytes.NewBuffer([]byte("test")))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOnSyncRequest_ValidRequest(t *testing.T) {
	router := gin.Default()
	broadcaster := &MockBroadcaster{}
	ginCSTV := gotv.NewGinCSTV(nil, nil, broadcaster)

	router.GET("/:token/sync", ginCSTV.OnSyncRequest())

	req, _ := http.NewRequest("GET", "/valid_token/sync", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var response gotv.Sync
	json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 200, response.Tick)
}

func TestOnSyncRequest_NotFound(t *testing.T) {
	router := gin.Default()
	broadcaster := &MockBroadcaster{}
	ginCSTV := gotv.NewGinCSTV(nil, nil, broadcaster)

	router.GET("/:token/sync", ginCSTV.OnSyncRequest())

	req, _ := http.NewRequest("GET", "/not_found/sync", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

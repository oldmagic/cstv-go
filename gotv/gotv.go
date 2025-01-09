package gotv

import (
	"io"
	"time"
)

// Auth handles authentication.
type Auth interface {
	Auth(token, auth string) error
}

// Store manages writing CSTV+ fragments.
type Store interface {
	OnStart(token string, fragment int, f StartFrame) error
	OnFull(token string, fragment int, tick int, at time.Time, r io.Reader) error
	OnDelta(token string, fragment int, endTick int, at time.Time, final bool, r io.Reader) error
}

// Broadcaster manages reading CSTV+ fragments.
type Broadcaster interface {
	GetSync(token string, fragment int) (Sync, error)
	GetSyncLatest(token string) (Sync, error)
	GetStart(token string, fragment int) (io.ReadCloser, error)
	GetFull(token string, fragment int) (io.ReadCloser, error)
	GetDelta(token string, fragment int) (io.ReadCloser, error)
}

// Fragment contains both full and delta data.
type Fragment struct {
	At      time.Time
	Tick    int
	Final   *bool
	EndTick int
	Full    []byte
	Delta   []byte
}

// StartFrame represents a starting fragment.
type StartFrame struct {
	At       time.Time
	TPS      float64
	Protocol int
	Map      string
	Body     []byte
}

// Sync represents sync JSON response.
type Sync struct {
	Tick             int     `json:"tick"`
	EndTick          int     `json:"endtick,omitempty"`
	RealTimeDelay    float64 `json:"rtdelay,omitempty"`
	ReceiveAge       float64 `json:"rcvage,omitempty"`
	Fragment         int     `json:"fragment"`
	SignupFragment   int     `json:"signup_fragment"`
	TPS              int     `json:"tps"`
	KeyframeInterval float64 `json:"keyframe_interval,omitempty"`
	TokenRedirect    string  `json:"token_redirect,omitempty"`
	Map              string  `json:"map"`
	Protocol         int     `json:"protocol"`
}

// StartQuery represents the query params for START requests.
type StartQuery struct {
	Tick int    `query:"tick"`
	TPS  int    `query:"tps"`
	Map  string `query:"map"`
}

// FullQuery represents the query params for FULL requests.
type FullQuery struct {
	Tick int `query:"tick"`
}

// DeltaQuery represents the query params for DELTA requests.
type DeltaQuery struct {
	EndTick int   `query:"endtick"`
	Final   *bool `query:"final"`
}

// SyncQuery represents the query params for SYNC requests.
type SyncQuery struct {
	Fragment *int `query:"fragment"`
}

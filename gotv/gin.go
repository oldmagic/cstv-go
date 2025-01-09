package gotv

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinCSTV manages the CSTV+ API handlers for the Gin framework.
type GinCSTV struct {
	auth        Auth
	store       Store
	broadcaster Broadcaster
}

// NewGinCSTV creates a new instance of GinCSTV.
func NewGinCSTV(auth Auth, store Store, broadcaster Broadcaster) *GinCSTV {
	return &GinCSTV{auth: auth, store: store, broadcaster: broadcaster}
}

// CheckAuthMiddleware verifies authentication via headers.
func (g *GinCSTV) CheckAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		authHeader := c.Request.Header.Get("X-Origin-Auth")

		if err := g.auth.Auth(token, authHeader); err != nil {
			logrus.Warnf("Unauthorized access attempt: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

// OnStartFragment handles the start fragment request.
func (g *GinCSTV) OnStartFragment() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fragment number"})
			return
		}

		var query StartQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}

		if err := g.store.OnFull(token, fragment, query.Tick, time.Now(), c.Request.Body); err != nil {
			if errors.Is(err, ErrMatchNotFound) {
				c.JSON(http.StatusResetContent, gin.H{"error": "Match not found"})
				return
			}
			logrus.Errorf("Failed to process start fragment: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}

// OnFullFragment handles the full fragment request.
func (g *GinCSTV) OnFullFragment() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fragment number"})
			return
		}

		var query FullQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}

		if err := g.store.OnFull(token, fragment, query.Tick, time.Now(), c.Request.Body); err != nil {
			if errors.Is(err, ErrMatchNotFound) {
				c.JSON(http.StatusResetContent, gin.H{"error": "Match not found"})
				return
			}
			if errors.Is(err, ErrFragmentNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Fragment not found"})
				return
			}
			logrus.Errorf("Failed to process full fragment: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}

// OnSyncRequest handles sync JSON requests.
func (g *GinCSTV) OnSyncRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		var query SyncQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}

		var syncData Sync
		var err error
		if query.Fragment != nil {
			syncData, err = g.broadcaster.GetSync(token, *query.Fragment)
		} else {
			syncData, err = g.broadcaster.GetSyncLatest(token)
		}

		if err != nil {
			logrus.Warnf("Sync request error: %v", err)
			if errors.Is(err, ErrMatchNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
			} else if errors.Is(err, ErrFragmentNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Fragment not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, syncData)
	}
}

// SetupStoreHandlers registers store-related routes.
func SetupStoreHandlers(g *GinCSTV, r *gin.RouterGroup) {
	r.Use(g.CheckAuthMiddleware())
	r.POST("/:token/:fragment_number/start", g.OnStartFragment())
	r.POST("/:token/:fragment_number/full", g.OnFullFragment())
}

// SetupBroadcasterHandlers registers broadcaster-related routes.
func SetupBroadcasterHandlers(g *GinCSTV, r *gin.RouterGroup) {
	r.Use(g.CheckAuthMiddleware())
	r.GET("/:token/sync", g.OnSyncRequest())
}

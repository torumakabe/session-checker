package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router with session management and routes.
// It supports both Redis and cookie-based session stores.
//
// Parameters:
//   - s: Redis server address (hostname:port). If empty, cookie-based session store is used.
//   - p: Redis password.
//
// Returns:
//   - *gin.Engine: The initialized Gin router.
func SetupRouter(s, p string) *gin.Engine {
	r := gin.Default()
	hostname, _ := os.Hostname()

	if s != "" {
		log.Printf("using redis for session store. server: %s\n", s)
		store, err := redis.NewStore(10, "tcp", s, p, []byte("secret"))
		if err != nil {
			log.Fatalf("Unable to connect to redis: %s\n", err)
		}
		r.Use(sessions.Sessions("mysession", store))
	} else {
		log.Println("using cookie for session store.")
		store := cookie.NewStore([]byte("secret"))
		store.Options(sessions.Options{
			SameSite: http.SameSiteStrictMode,
		})
		r.Use(sessions.Sessions("mysession", store))
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi, GET /incr to check session.")
	})

	r.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}

		session.Set("count", count)
		err := session.Save()
		if err != nil {
			log.Printf("session save error: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to save session",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"count":    count,
			"hostname": hostname,
		})
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "OK",
			"hostname": hostname,
		})
	})

	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "OK",
			"hostname": hostname,
		})
	})

	return r
}

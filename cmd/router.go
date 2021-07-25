package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func setupRouter(s, p string) *gin.Engine {
	r := gin.Default()

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
		session.Save()
		hostname, _ := os.Hostname()
		c.JSON(200, gin.H{
			"count":    count,
			"hostname": hostname,
		})
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	return r
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var redisServer string

func main() {
	flag.StringVarP(&redisServer, "redis-server", "r", "", "set redis server hostname:port when you use")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
	viper.SetEnvPrefix("SESSION_CHECKER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	r := gin.Default()

	if viper.GetString("redis-server") != "" {
		log.Printf("using redis for session store. server: %s\n", viper.GetString("redis-server"))
		store, err := redis.NewStore(10, "tcp", viper.GetString("redis-server"), "", []byte("secret"))
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		log.Println("got signal. shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	log.Printf("start receiving at %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

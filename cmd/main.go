package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/torumakabe/session-checker/router"
)

var (
	redisServer   string
	redisPassword string
)

func main() {
	pflag.StringVarP(&redisServer, "redis-server", "r", "", "set redis server hostname:port when you use")
	pflag.StringVarP(&redisPassword, "redis-password", "p", "", "set redis password when you use")

	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("flag parse error: %s", err)
	}
	viper.SetEnvPrefix("SESSION_CHECKER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	r := router.SetupRouter(viper.GetString("redis-server"), viper.GetString("redis-password"))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		log.Println("got signal. shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Fatalf("server shutdown error: %s", err)
		}
	}()
	log.Printf("start receiving at %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

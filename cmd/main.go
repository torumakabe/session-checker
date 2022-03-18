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

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	redisServer   string
	redisPassword string
)

func main() {
	flag.StringVarP(&redisServer, "redis-server", "r", "", "set redis server hostname:port when you use")
	flag.StringVarP(&redisPassword, "redis-password", "p", "", "set redis password when you use")

	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
	viper.SetEnvPrefix("SESSION_CHECKER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	r := setupRouter(viper.GetString("redis-server"), viper.GetString("redis-password"))

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

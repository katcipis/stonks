package main

import (
	"context"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/katcipis/stonks/api"
	"github.com/katcipis/stonks/auth"
	"github.com/katcipis/stonks/users/manager"
	"github.com/katcipis/stonks/users/storage"
)

type Config struct {
	UsersDBHost     string
	UsersDBName     string
	UsersDBUser     string
	UsersDBPassword string
}

func main() {

	cfg := loadCfg()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	usersStorage, err := storage.New(
		ctx,
		cfg.UsersDBHost,
		cfg.UsersDBName,
		cfg.UsersDBUser,
		cfg.UsersDBPassword,
	)
	if err != nil {
		panic(err)
	}

	authorizer := auth.New()
	usersManager := manager.New(authorizer, usersStorage)

	service := api.New(usersManager, api.Config{
		CreateUserTimeout: 10 * time.Second,
	})

	// Usually I add a port flag parameter, running against time :-)
	// Server wide timeouts can be dangerous if you have some endpoint
	// that does streaming of results (or long polling). Since this is
	// not the case I set a timeout to ensure the server never hangs.
	server := &http.Server{
		Addr:           ":8080",
		Handler:        service,
		ReadTimeout:    time.Minute,
		WriteTimeout:   time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	log.Infof("running users manager service at %q", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func loadCfg() Config {
	return Config{
		UsersDBHost:     loadenv("USERS_DB_HOST", "usersdb"),
		UsersDBName:     loadenv("USERS_DB_NAME", "testing"),
		UsersDBUser:     loadenv("USERS_DB_USER", "testing"),
		UsersDBPassword: loadenv("USERS_DB_PASSWORD", "testing"),
	}
}

func loadenv(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}

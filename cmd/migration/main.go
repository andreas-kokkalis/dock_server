package main

import (
	"errors"
	"log"
	"regexp"

	"github.com/Davmuz/gqt"
	"github.com/andreas-kokkalis/dock-server/pkg/api/store"
	"github.com/andreas-kokkalis/dock-server/pkg/config"
	"github.com/caarlos0/env"
)

type envVars struct {
	Mode string `env:"MODE"`
}

var validMode = regexp.MustCompile(`^(local)`)

var errInvalidMode = errors.New("Invalid environment variable MODE\n Allowed values [local]")

var configDir = "./conf"

func main() {
	vars := envVars{}
	err := env.Parse(&vars)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if !validMode.MatchString(vars.Mode) {
		log.Fatal(errInvalidMode)
	}

	// Initialize the configuration manager
	var c config.Config
	c, err = config.NewConfig(configDir, vars.Mode)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Postgres storage
	var db *store.DB
	db, err = store.NewDB("postgres", c.GetPGConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	migrateData(db)
}

func migrateData(db *store.DB) {
	// Create database schema
	_ = gqt.Add("templates/pgsql", "*.pgsql")
	_, err := db.Query(gqt.Get("createSchema"))
	if err != nil {
		log.Fatal(err)
	}
	// Insert Data
	_, err = db.Query(gqt.Get("migrateData"))
	if err != nil {
		log.Fatal(err)
	}
}

package srv

import (
	"log"
	"os"

	"github.com/andreas-kokkalis/dock-server/conf"
	"github.com/andreas-kokkalis/dock-server/db"
	"github.com/andreas-kokkalis/dock-server/dc"
)

// InitDep initializes connections with the viper, docker api client, portMapper, redis and postgres
func InitDep() {
	// Load static configuration strings from conf/conf.yaml
	err := conf.InitConf("./conf")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the Docker API Client
	dc.APIClientInit(conf.GetVal("dc.docker.api.version"), conf.GetVal("dc.docker.api.host"))

	// Initialize the port mappings
	dc.ContainerPortsInitialize(200)

	// Initialize Redis storage
	db.InitRedisClient()

	// Initialize Postgres storage
	db.InitPostgres()

	// Create Schema and insert data if mode is set to dev
	mode := os.Getenv("MODE")
	if mode == "dev" {
		db.MigrateData()
	}

	// Start a goroute that will run the PeriodicChecker
	go dc.PeriodicChecker()

}

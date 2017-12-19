package schema

import (
	"regexp"

	"github.com/Davmuz/gqt"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
	"github.com/andreas-kokkalis/dock_server/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO: Add logging

var (
	// ConfigDir flag indicates the location of the conf.yaml file
	ConfigDir string

	// Env flag indicates the environment that the server will run.
	Env           string
	vEnv          = regexp.MustCompile(`^(local)`)
	errInvalidEnv = errors.New("Allowed env values are [local]")
)

// Create command starts the HTTP API server
var Create = func(cmd *cobra.Command, args []string) (err error) {
	if !vEnv.MatchString(Env) {
		return errInvalidEnv
	}

	// Initialize the configuration manager
	var c *config.Config
	if c, err = config.NewConfig(ConfigDir, Env); err != nil {
		return err
	}

	// Initialize Postgres storage
	var db *store.DB
	if db, err = store.NewDB(c.GetPGConnectionString()); err != nil {
		return errors.Wrap(err, "Unable to connect to the database")
	}

	return migrateData(db)
}

func migrateData(db *store.DB) error {
	// Create database schema
	_ = gqt.Add("templates/pgsql", "*.pgsql")
	if _, err := db.Query(gqt.Get("createSchema")); err != nil {
		return err
	}
	// Insert Data
	if _, err := db.Query(gqt.Get("migrateData")); err != nil {
		return err
	}
	return nil
}

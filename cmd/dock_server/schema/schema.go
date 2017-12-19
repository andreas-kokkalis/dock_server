package schema

import (
	"regexp"

	"github.com/andreas-kokkalis/dock_server/pkg/config"
	"github.com/andreas-kokkalis/dock_server/pkg/util/dbutil"
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

	// ScriptDir flag indicates the path where the database scripts are located.
	ScriptDir string
)

func initConf() (c *config.Config, err error) {
	if !vEnv.MatchString(Env) {
		return nil, errInvalidEnv
	}

	// Initialize the configuration manager
	if c, err = config.NewConfig(ConfigDir, Env); err != nil {
		return nil, err
	}
	return c, nil
}

// Create command creates the database schema
var Create = func(cmd *cobra.Command, args []string) (err error) {

	c, err := initConf()
	if err != nil {
		return err
	}
	// Initialize Postgres storage
	dbm, err := dbutil.NewDBManager(c.GetPGConnectionString(), ScriptDir)
	if err != nil {
		return errors.Wrap(err, "Unable to connect to the database")
	}
	return dbm.CreateSchema()
}

// Drop command drops the database schema
var Drop = func(cmd *cobra.Command, args []string) (err error) {

	c, err := initConf()
	if err != nil {
		return err
	}
	// Initialize Postgres storage
	dbm, err := dbutil.NewDBManager(c.GetPGConnectionString(), ScriptDir)
	if err != nil {
		return errors.Wrap(err, "Unable to connect to the database")
	}
	return dbm.DropSchema()
}

// Insert command inserts basic required data for in the database schema
var Insert = func(cmd *cobra.Command, args []string) (err error) {

	c, err := initConf()
	if err != nil {
		return err
	}
	// Initialize Postgres storage
	dbm, err := dbutil.NewDBManager(c.GetPGConnectionString(), ScriptDir)
	if err != nil {
		return errors.Wrap(err, "Unable to connect to the database")
	}
	return dbm.InsertSchema()

}

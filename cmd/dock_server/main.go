package main

import (
	"log"

	"github.com/andreas-kokkalis/dock_server/cmd/dock_server/schema"
	"github.com/andreas-kokkalis/dock_server/cmd/dock_server/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dock",
	Short: "Command line client for interacting the the dock API server",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "wrapper command to interact with the dock API server",
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the api server",
	RunE:  server.Start,
	Args:  cobra.NoArgs,
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "wrapper command to interact with the database for integration purposes",
}

var schemaCreate = &cobra.Command{
	Use:   "create",
	Short: "create the database schema and migrates base data",
	RunE:  schema.Create,
	Args:  cobra.NoArgs,
}

var schemaDrop = &cobra.Command{
	Use:   "drop",
	Short: "drop the database schema",
	RunE:  schema.Drop,
	Args:  cobra.NoArgs,
}

var schemaInsert = &cobra.Command{
	Use:   "insert",
	Short: "insert data to the database schema",
	RunE:  schema.Insert,
	Args:  cobra.NoArgs,
}

func main() {

	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVarP(&server.ConfigDir, "conf", "c", "./conf", "The directory where the conf.yaml file is located.")
	serverCmd.PersistentFlags().StringVarP(&server.Env, "env", "e", "local", "The environment target for the api.")
	serverCmd.AddCommand(serverStartCmd)

	rootCmd.AddCommand(schemaCmd)
	schemaCmd.PersistentFlags().StringVarP(&schema.ConfigDir, "conf", "c", "./conf", "The directory where the conf.yaml file is located.")
	schemaCmd.PersistentFlags().StringVarP(&schema.ScriptDir, "script", "s", "./scripts/db", "The directory where the database scripts are located.")
	schemaCmd.PersistentFlags().StringVarP(&schema.Env, "env", "e", "local", "The environment target for the configuration")
	schemaCmd.AddCommand(schemaCreate)
	schemaCmd.AddCommand(schemaDrop)
	schemaCmd.AddCommand(schemaInsert)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

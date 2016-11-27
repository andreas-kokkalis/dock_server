package main

import (
	"log"
	"net/http"
	"os"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route"
	"github.com/andreas-kokkalis/dock-server/srv"
	"github.com/julienschmidt/httprouter"
)

func main() {

	dc.ClientInit("", "")
	srv.InitPortMappings(100)
	srv.InitRedisClient()
	srv.InitPostgres()

	// Create Schema and insert data if mode is set to dev
	mode := os.Getenv("MODE")
	if mode == "dev" {
		srv.MigrateData()
	}

	router := httprouter.New()

	// List of Routes

	// Containers
	router.GET("/v0/containers", route.GetContainers)
	router.GET("/v0/containers/:status", route.GetContainers)
	router.POST("/v0/containers/run", route.RunContainer)
	router.POST("/v0/containers/commit", route.CommitContainer)
	router.DELETE("/v0/containers/kill/:id", route.KillContainer)

	// Images
	router.GET("/v0/images", route.ListImages)
	router.GET("/v0/images/history/:id", route.GetImageHistory)
	router.DELETE("/v0/images/:id", route.RemoveImage)

	// Admin
	router.POST("/v0/login/", route.AdminLogin)

	// Start the server
	err := http.ListenAndServeTLS(":8080", "ssl/server.pem", "ssl/server.key", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

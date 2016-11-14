package main

import (
	"log"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route"
	"github.com/andreas-kokkalis/dock-server/srv"
	"github.com/julienschmidt/httprouter"
)

func main() {

	dc.ClientInit("", "")
	srv.InitPortMappings(100)
	srv.InitRedisClient()
	router := httprouter.New()

	// List of Routes
	router.GET("/v0/containers", route.GetContainers)
	router.GET("/v0/containers/:status", route.GetContainers)
	router.GET("/v0/images", route.ListImages)
	router.POST("/v0/run", route.RunContainer)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

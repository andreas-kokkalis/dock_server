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

	// Containers
	router.GET("/v0/containers", route.GetContainers)
	router.GET("/v0/containers/:status", route.GetContainers)
	router.POST("/v0/run", route.RunContainer)

	// Images
	router.GET("/v0/images", route.ListImages)
	router.GET("/v0/images/history/:id", route.GetImageHistory)
	router.DELETE("/v0/images/:id", route.RemoveImage)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

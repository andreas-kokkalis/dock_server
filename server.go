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

	/****************
	* ADMIN ROUTES
	****************/
	// Login to Panel
	router.GET("/v0/admin/logout", route.AdminLogout)
	router.POST("/v0/admin/login", route.AdminLogin)
	// Container actions
	router.GET("/v0/admin/containers/list", route.AuthAdmin(route.GetContainers))
	router.GET("/v0/admin/containers/list/:status", route.AuthAdmin(route.GetContainers))
	router.POST("/v0/admin/containers/run/:id", route.AuthAdmin(route.AdminRunContainer))
	router.POST("/v0/admin/containers/commit/:id", route.AuthAdmin(route.CommitContainer))
	router.POST("/v0/admin/containers/kill/:id", route.AuthAdmin(route.AdminKillContainer))
	// Image actions
	router.GET("/v0/admin/images", route.AuthAdmin(route.ListImages))
	router.GET("/v0/admin/images/history/:id", route.AuthAdmin(route.GetImageHistory))
	router.DELETE("/v0/admin/images/delete/:id", route.AuthAdmin(route.RemoveImage))

	/****************
	* USER ROUTES
	****************/
	// LTILaunch	- id is the imageID
	router.POST("/v0/lti/launch/:id", route.OAuth(route.LTILaunch))

	/****************
	* ADMIN FRONTEND
	****************/
	// Serve the frontend files for the admin panel
	router.ServeFiles("/ui/*filepath", http.Dir("./public/"))

	// Start the server
	err := http.ListenAndServeTLS(":8080", "ssl/server.pem", "ssl/server.key", router)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

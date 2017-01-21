package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/andreas-kokkalis/dock-server/conf"
	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route"
	"github.com/andreas-kokkalis/dock-server/srv"
	"github.com/julienschmidt/httprouter"
)

func main() {
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
	srv.InitRedisClient()

	// Initialize Postgres storage
	srv.InitPostgres()

	// Create Schema and insert data if mode is set to dev
	mode := os.Getenv("MODE")
	if mode == "dev" {
		srv.MigrateData()
	}

	go dc.PeriodicChecker()

	// Initialize the  httprouter
	router := httprouter.New()

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
	router.DELETE("/v0/admin/containers/kill/:id", route.AdminKillContainer)
	// Image actions
	router.GET("/v0/admin/images", route.AuthAdmin(route.ListImages))
	router.GET("/v0/admin/images/history/:id", route.AuthAdmin(route.GetImageHistory))
	router.DELETE("/v0/admin/images/delete/:id", route.AuthAdmin(route.RemoveImage))
	router.GET("/v0/admin/ping", route.AuthAdmin(route.Ping))

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
	myServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      router,
	}
	err = myServer.ListenAndServeTLS("conf/ssl/server.pem", "conf/ssl/server.key")

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

package main

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/andreas-kokkalis/dock-server/pkg/api/auth"
	"github.com/andreas-kokkalis/dock-server/pkg/api/docker"
	"github.com/andreas-kokkalis/dock-server/pkg/api/image"
	"github.com/andreas-kokkalis/dock-server/pkg/api/store"
	"github.com/andreas-kokkalis/dock-server/pkg/config"
	"github.com/caarlos0/env"
	"github.com/julienschmidt/httprouter"
)

type envVars struct {
	Mode string `env:"MODE"`
}

var validMode = regexp.MustCompile(`^(local)`)

var errInvalidMode = errors.New("Invalid environment variable MODE\n Allowed values [local]")

var configDir = "./conf"

// TODO: figure out what to do when migration is required

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

	// Initialize Redis storage
	var redis *store.Redis
	redis, err = store.NewRedisClient(c.GetRedisConfig())
	if err != nil {
		log.Fatal(err)
	}
	// Initialize Redis repository
	redisRepository := store.NewRedisRepo(redis)

	// Initialize Docker Remote API Client
	var dockerClient *docker.DockerCli
	dockerClient, err = docker.NewAPIClient(c.GetDockerConfig())
	if err != nil {
		log.Fatal(err)
	}
	// Initialize docker repository
	dockerRepository := docker.NewRepo(dockerClient, c.GetDockerConfig())

	// Initialize PortMapper
	mapper := docker.NewPortMapper(redisRepository, c.GetAPIPorts())
	// Start a goroute that will run the PeriodicChecker
	go docker.PeriodicChecker(dockerRepository, mapper, redisRepository)

	// Initialize the  httprouter
	router := httprouter.New()

	// Auth Service
	authService := auth.NewService(db, redisRepository, dockerRepository)
	router.GET("/v0/admin/logout", auth.AdminLogout(authService))
	router.POST("/v0/admin/login", auth.AdminLogin(authService))

	// Image service
	imageService := image.NewService(db, redisRepository, dockerRepository)
	router.GET("/v0/admin/images", auth.AuthAdmin(authService, image.ListImages(imageService)))
	router.GET("/v0/admin/images/history/:id", auth.AuthAdmin(authService, image.GetImageHistory(imageService)))
	router.DELETE("/v0/admin/images/delete/:id", auth.AuthAdmin(authService, image.RemoveImage(imageService)))

	/****************
	* ADMIN ROUTES
	****************/
	// // Container actions
	// router.GET("/v0/admin/containers/list", route.AuthAdmin(route.GetContainers))
	// router.GET("/v0/admin/containers/list/:status", route.AuthAdmin(route.GetContainers))
	// router.POST("/v0/admin/containers/run/:id", route.AuthAdmin(route.AdminRunContainer))
	// router.POST("/v0/admin/containers/commit/:id", route.AuthAdmin(route.CommitContainer))
	// router.DELETE("/v0/admin/containers/kill/:id", route.AdminKillContainer)
	// // Image actions
	/****************
	* USER ROUTES
	****************/
	// LTILaunch	- id is the imageID
	// router.POST("/v0/lti/launch/:id", route.OAuth(route.LTILaunch))

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

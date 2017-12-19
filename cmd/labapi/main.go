package main

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/andreas-kokkalis/dock_server/pkg/api/auth"
	"github.com/andreas-kokkalis/dock_server/pkg/api/container"
	"github.com/andreas-kokkalis/dock_server/pkg/api/docker"
	"github.com/andreas-kokkalis/dock_server/pkg/api/image"
	"github.com/andreas-kokkalis/dock_server/pkg/api/lti"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
	"github.com/andreas-kokkalis/dock_server/pkg/config"
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
	var c *config.Config
	c, err = config.NewConfig(configDir, vars.Mode)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Postgres storage
	var db *store.DB
	db, err = store.NewDB(c.GetPGConnectionString())
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

	// Initialize PortMapper
	mapper := docker.NewPortMapper(redisRepository, c.GetAPIPorts())
	// Initialize Docker Remote API Client

	var dockerClient *docker.APIClient
	dockerClient, err = docker.NewAPIClient(c.GetDockerConfig())
	if err != nil {
		log.Fatal(err)
	}
	// Initialize docker repository
	dockerRepository := docker.NewRepo(dockerClient, c.GetDockerConfig())

	// Start a goroute that will run the PeriodicChecker
	go docker.PeriodicChecker(dockerRepository, mapper, redisRepository)

	// Initialize the  httprouter
	router := httprouter.New()

	// Auth Service
	authService := auth.NewService(db, redisRepository)
	router.GET("/v0/admin/logout", auth.AdminLogout(authService))
	router.POST("/v0/admin/login", auth.AdminLogin(authService))

	// Image service
	imageService := image.NewService(db, redisRepository, dockerRepository)
	router.GET("/v0/admin/images", auth.SessionAuth(authService, image.ListImages(imageService)))
	router.GET("/v0/admin/images/history/:id", auth.SessionAuth(authService, image.GetImageHistory(imageService)))
	router.DELETE("/v0/admin/images/delete/:id", auth.SessionAuth(authService, image.RemoveImage(imageService)))

	// Container service
	containerService := container.NewService(db, redisRepository, dockerRepository, mapper)
	router.POST("/v0/admin/containers/run/:id", auth.SessionAuth(authService, container.AdminRunContainer(containerService)))
	router.DELETE("/v0/admin/containers/kill/:id", auth.SessionAuth(authService, container.AdminKillContainer(containerService)))
	router.POST("/v0/admin/containers/commit/:id", auth.SessionAuth(authService, container.CommitContainer(containerService)))
	router.GET("/v0/admin/containers/list", auth.SessionAuth(authService, container.GetContainers(containerService)))
	router.GET("/v0/admin/containers/list/:status", auth.SessionAuth(authService, container.GetContainers(containerService)))

	// LTI service
	ltiService := lti.NewService(db, redisRepository, dockerRepository, mapper)
	router.POST("/v0/lti/launch/:id", auth.OAuth(authService, lti.Launch(ltiService)))

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

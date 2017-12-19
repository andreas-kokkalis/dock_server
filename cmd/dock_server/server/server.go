package server

import (
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/andreas-kokkalis/dock_server/pkg/api/auth"
	"github.com/andreas-kokkalis/dock_server/pkg/api/container"
	"github.com/andreas-kokkalis/dock_server/pkg/api/docker"
	"github.com/andreas-kokkalis/dock_server/pkg/api/image"
	"github.com/andreas-kokkalis/dock_server/pkg/api/lti"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
	"github.com/andreas-kokkalis/dock_server/pkg/config"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/db"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/redis"
	"github.com/julienschmidt/httprouter"
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

// Start command starts the HTTP API server
var Start = func(cmd *cobra.Command, args []string) (err error) {

	if !vEnv.MatchString(Env) {
		return errInvalidEnv
	}

	// Initialize the configuration manager
	var c *config.Config
	if c, err = config.NewConfig(ConfigDir, Env); err != nil {
		return err
	}

	// Initialize Postgres storage
	var dbConn *db.DB
	if dbConn, err = db.NewDB(c.GetPGConnectionString()); err != nil {
		return errors.Wrap(err, "Unable to connect to the database")
	}

	// Initialize Redis storage
	var redisCli redis.Redis
	if redisCli, err = redis.NewClient(c.GetRedisConfig()); err != nil {
		return errors.Wrap(err, "Unable to connect to redis")
	}
	// Initialize Redis repository
	redisRepository := store.NewRedisRepo(redisCli)

	// Initialize PortMapper
	mapper := docker.NewPortMapper(redisRepository, c.GetAPIPorts())
	// Initialize Docker Remote API Client

	var dockerClient *docker.APIClient
	if dockerClient, err = docker.NewAPIClient(c.GetDockerConfig()); err != nil {
		return err
	}
	// Initialize docker repository
	dockerRepository := docker.NewRepo(dockerClient, c.GetDockerConfig())

	// Start a goroute that will run the PeriodicChecker
	go docker.PeriodicChecker(dockerRepository, mapper, redisRepository)

	// Initialize the  httprouter
	router := httprouter.New()

	// Auth Service
	authService := auth.NewService(dbConn, redisRepository)
	router.GET("/v0/admin/logout", auth.AdminLogout(authService))
	router.POST("/v0/admin/login", auth.AdminLogin(authService))

	// Image service
	imageService := image.NewService(dbConn, redisRepository, dockerRepository)
	router.GET("/v0/admin/images", auth.SessionAuth(authService, image.ListImages(imageService)))
	router.GET("/v0/admin/images/history/:id", auth.SessionAuth(authService, image.GetImageHistory(imageService)))
	router.DELETE("/v0/admin/images/delete/:id", auth.SessionAuth(authService, image.RemoveImage(imageService)))

	// Container service
	containerService := container.NewService(dbConn, redisRepository, dockerRepository, mapper)
	router.POST("/v0/admin/containers/run/:id", auth.SessionAuth(authService, container.AdminRunContainer(containerService)))
	router.DELETE("/v0/admin/containers/kill/:id", auth.SessionAuth(authService, container.AdminKillContainer(containerService)))
	router.POST("/v0/admin/containers/commit/:id", auth.SessionAuth(authService, container.CommitContainer(containerService)))
	router.GET("/v0/admin/containers/list", auth.SessionAuth(authService, container.GetContainers(containerService)))
	router.GET("/v0/admin/containers/list/:status", auth.SessionAuth(authService, container.GetContainers(containerService)))

	// LTI service
	ltiService := lti.NewService(dbConn, redisRepository, dockerRepository, mapper)
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
		Addr:         c.GetAPIServerPort(),
		Handler:      router,
	}
	err = myServer.ListenAndServeTLS("conf/ssl/server.pem", "conf/ssl/server.key")
	if err != nil {
		return errors.Wrap(err, "ListenAndServe")
	}
	return nil
}

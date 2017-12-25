package server

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/andreas-kokkalis/dock_server/pkg/api/auth"
	"github.com/andreas-kokkalis/dock_server/pkg/api/container"
	"github.com/andreas-kokkalis/dock_server/pkg/api/image"
	"github.com/andreas-kokkalis/dock_server/pkg/api/lti"
	"github.com/andreas-kokkalis/dock_server/pkg/api/portmapper"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories"
	"github.com/andreas-kokkalis/dock_server/pkg/config"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/docker"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres"
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
	var dbConn *postgres.DB
	if dbConn, err = postgres.NewDB(c.GetPGConnectionString()); err != nil {
		return errors.Wrap(err, "Unable to connect to the database")
	}

	// Initialize Redis storage
	var redisCli redis.Redis
	if redisCli, err = redis.NewClient(c.GetRedisConfig()); err != nil {
		return errors.Wrap(err, "Unable to connect to redis")
	}
	// Initialize Redis repository
	redisRepository := repositories.NewRedisRepo(redisCli)

	// Initialize PortMapper
	mapper := portmapper.NewPortMapper(redisRepository, c.GetAPIPorts())
	// Initialize Docker Remote API Client

	var dockerClient *docker.APIClient
	if dockerClient, err = docker.NewAPIClient(c.GetDockerConfig()); err != nil {
		return err
	}
	// Initialize docker repository
	dockerRepository := repositories.NewDockerRepository(dockerClient, c.GetDockerConfig())

	// Start a goroute that will run the PeriodicChecker
	go portmapper.PeriodicChecker(dockerRepository, mapper, redisRepository)

	// Initialize the  httprouter
	router := httprouter.New()

	// Auth Service
	adminRepo := repositories.NewAdminDBRepository(dbConn)
	authService := auth.NewService(adminRepo, redisRepository)
	router.GET("/v0/admin/logout", authService.AdminLogout)
	router.POST("/v0/admin/login", authService.AdminLogin)

	// Image service
	imgService := image.NewService(redisRepository, dockerRepository)
	router.GET("/v0/admin/images", authService.SessionAuth(imgService.ListImages))
	router.GET("/v0/admin/images/history/:id", authService.SessionAuth(imgService.GetImageHistory))
	router.DELETE("/v0/admin/images/delete/:id", authService.SessionAuth(imgService.RemoveImage))

	// Container service
	cntService := container.NewService(redisRepository, dockerRepository, mapper)
	router.POST("/v0/admin/containers/run/:id", authService.SessionAuth(cntService.AdminRunContainer))
	router.DELETE("/v0/admin/containers/kill/:id", authService.SessionAuth(cntService.AdminKillContainer))
	router.POST("/v0/admin/containers/commit/:id", authService.SessionAuth(cntService.CommitContainer))
	router.GET("/v0/admin/containers/list", authService.SessionAuth(cntService.GetContainers))
	router.GET("/v0/admin/containers/list/:status", authService.SessionAuth(cntService.GetContainers))

	// LTI service
	ltiService := lti.NewService(dbConn, redisRepository, dockerRepository, mapper)
	router.POST("/v0/lti/launch/:id", authService.OAuth(lti.Launch(ltiService)))

	/****************
	* ADMIN FRONTEND
	****************/
	// Serve the frontend files for the admin panel
	router.ServeFiles("/ui/*filepath", http.Dir("./public/"))

	// Start the server
	myServer := &http.Server{
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 10 * time.Second,
		Addr:    c.GetAPIServerPort(),
		Handler: logger{router},
	}
	err = myServer.ListenAndServeTLS("conf/ssl/server.pem", "conf/ssl/server.key")
	if err != nil {
		return errors.Wrap(err, "ListenAndServe")
	}
	return nil
}

type logger struct {
	handler http.Handler
}

func (l logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		log.Printf("[took: %v] %s %s", time.Since(start), r.Method, r.URL.Path)
	}()
	l.handler.ServeHTTP(w, r)
}

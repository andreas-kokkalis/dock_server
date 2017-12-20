package integration

import (
	"encoding/json"
	"log"
	"net/http"
	"path"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/docker"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
	"github.com/andreas-kokkalis/dock_server/pkg/config"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/redis"
	"github.com/andreas-kokkalis/dock_server/pkg/util/dbutil"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	environment = "local"
	confDir     = "conf"
	scriptDir   = "scripts/db"
)

// Spec struct is used to hold information between different suites of integration tests
type Spec struct {
	// Top level directory
	TopDir string

	// Configuration
	Config *config.Config

	// Postgres
	DBManager *dbutil.DBManager

	// Redis
	Redis     redis.Redis
	RedisRepo *store.RedisRepo

	// Docker
	DockerCLI  *docker.APIClient
	DockerRepo *docker.Repo

	// Logger
	Log *log.Logger

	// Handler
	Handler http.Handler
}

// NewSpec initializes a spec struct with the given values
func NewSpec(topDir string) *Spec {
	return &Spec{
		TopDir: topDir,
		Log:    log.New(ginkgo.GinkgoWriter, "", log.LstdFlags),
	}
}

// InitConfig returns a func to be used by BeforeSuite
// initializes configuration
func (s *Spec) InitConfig() func() {
	return func() {
		c, err := config.NewConfig(path.Join(s.TopDir, confDir), environment)
		gomega.Expect(err).To(gomega.BeNil(), "Init config")
		s.Config = c
	}
}

// InitDBConnection returns a func to be used by BeforeSuite
// establishes postgres connection
func (s *Spec) InitDBConnection() func() {
	return func() {
		db, err := dbutil.NewDBManager(s.Config.GetPGConnectionString(), path.Join(s.TopDir, scriptDir))
		gomega.Expect(err).To(gomega.BeNil(), "Connect Postgres")
		s.DBManager = db
	}
}

// CloseDBConnection returns a function that closes the Postgres connection poo;
func (s *Spec) CloseDBConnection() func() {
	return func() {
		gomega.Expect(s.DBManager.DB.Conn.Close()).To(gomega.BeNil(), "Disconnect Postgres")
	}
}

// RestoreDB drops the database schema, recreates it, and migrates data
func (s *Spec) RestoreDB() func() {
	return func() {

		err := s.DBManager.DropSchema()
		gomega.Expect(err).To(gomega.BeNil(), "dropping database tables")

		err = s.DBManager.CreateSchema()
		gomega.Expect(err).To(gomega.BeNil(), "creating database tables")

		err = s.DBManager.InsertSchema()
		gomega.Expect(err).To(gomega.BeNil(), "migrating data")
	}
}

// InitRedisConnection returns a func to be used by BeforeSuite
// establishes redis connection
func (s *Spec) InitRedisConnection() func() {
	return func() {
		redis, err := redis.NewClient(s.Config.GetRedisConfig())
		gomega.Expect(err).To(gomega.BeNil(), "Connect Redis")
		s.Redis = redis
		s.RedisRepo = store.NewRedisRepo(s.Redis)
	}
}

// CloseRedisConnection closes the connection pool to Redis
func (s *Spec) CloseRedisConnection() func() {
	return func() {
		gomega.Expect(s.Redis.Close()).To(gomega.BeNil(), "Disconnect Redis")
	}
}

// InitDockerRepo initializes the docker repository and the connection to docker API client
func (s *Spec) InitDockerRepo() func() {
	return func() {
		dockerClient, err := docker.NewAPIClient(s.Config.GetDockerConfig())
		gomega.Expect(err).To(gomega.BeNil(), "Init docker api client")
		s.DockerCLI = dockerClient
		s.DockerRepo = docker.NewRepo(dockerClient, s.Config.GetDockerConfig())
	}
}

// AssertAPICall performs an HTTP request, records the output and asserts if it matches against the expected response code and body.
func (s *Spec) AssertAPICall(request *Request, response *Response) {

	// Perform HTTP Request
	s.Handler.ServeHTTP(response.recorder, request.HTTPRequest)

	// Log request and response to stdout
	s.Log.Printf("\n------------------\n%s\n------------------\n", request.pretty())
	s.Log.Printf("\n------------------\n%s\n------------------\n", response.pretty())

	// Perform assertions
	gomega.Expect(response.Code()).To(gomega.Equal(response.expectedCode), "status codes do not match")

	var actualResponse api.Response
	err := json.Unmarshal(response.recorder.Body.Bytes(), &actualResponse)
	gomega.Expect(err).To(gomega.BeNil())

	diff, err := CompareRegexJSON(response.expectedBody, response.ToString(), s.TopDir)
	gomega.Expect(err).To(gomega.BeNil(), "Diff tool returned error")
	gomega.Expect(diff).To(gomega.Equal(""), "Diff is not empty")
}

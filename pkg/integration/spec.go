package integration

import (
	"log"
	"path"

	"github.com/Davmuz/gqt"
	"github.com/andreas-kokkalis/dock_server/pkg/api/docker"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
	"github.com/andreas-kokkalis/dock_server/pkg/config"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	environment = "local"
	confDir     = "conf"
	// TestDataDir is hardconfigured to be named testdata within each spec directory
	TestDataDir = "testdata"
)

// Spec struct is used to hold information between different suites of integration tests
type Spec struct {
	// Top level directory
	TopDir string

	// Configuration
	Config *config.Config

	// Postgres
	DB *store.DB

	// Redis
	Redis     *store.Redis
	RedisRepo *store.RedisRepo

	// Docker
	DockerCLI  *docker.APIClient
	DockerRepo *docker.Repo

	// Logger
	Log *log.Logger
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
		db, err := store.NewDB(s.Config.GetPGConnectionString())
		gomega.Expect(err).To(gomega.BeNil(), "Connect Postgres")
		s.DB = db
	}
}

// CloseDBConnection returns a function that closes the Postgres connection poo;
func (s *Spec) CloseDBConnection() func() {
	return func() {
		gomega.Expect(s.DB.Conn.Close()).To(gomega.BeNil(), "Disconnect Postgres")
	}
}

// RestoreDB drops the database schema and recreates it
func (s *Spec) RestoreDB() func() {
	return func() {

		err := gqt.Add(path.Join(s.TopDir, "templates/pgsql"), "*.pgsql")
		gomega.Expect(err).To(gomega.BeNil(), "load SQL template directory")

		_, err = s.DB.Query(gqt.Get("dropSchema"))
		gomega.Expect(err).To(gomega.BeNil(), "dropping database tables")

		_, err = s.DB.Query(gqt.Get("createSchema"))
		gomega.Expect(err).To(gomega.BeNil(), "creating database tables")

		// Insert Data
		_, err = s.DB.Query(gqt.Get("migrateData"))
		gomega.Expect(err).To(gomega.BeNil(), "migrating data")
	}
}

// InitRedisConnection returns a func to be used by BeforeSuite
// establishes redis connection
func (s *Spec) InitRedisConnection() func() {
	return func() {
		redis, err := store.NewRedisClient(s.Config.GetRedisConfig())
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

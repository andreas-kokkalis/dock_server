package container

import (
	"github.com/andreas-kokkalis/dock-server/pkg/api/docker"
	"github.com/andreas-kokkalis/dock-server/pkg/api/store"
)

// Service for image
type Service struct {
	db     *store.DB
	redis  *store.RedisRepo
	docker *docker.Repo
}

// NewService creates a new Image Service
func NewService(db *store.DB, redis *store.RedisRepo, docker *docker.Repo) Service {
	return Service{db, redis, docker}
}

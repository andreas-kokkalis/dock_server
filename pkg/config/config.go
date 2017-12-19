package config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	redis "gopkg.in/redis.v5"
)

// Config the configuration struct holding a pointer to Viper  configuration manager
type Config struct {
	viper *viper.Viper
	env   string
}

// NewConfig initializes a viper configuration and returns
func NewConfig(configDir string, env string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(configDir)
	v.SetConfigType("yaml")
	v.SetConfigName("conf")
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Could not read config file conf.yaml at location"+configDir)
	}
	c := &Config{
		env:   env,
		viper: v.Sub(env),
	}
	return c, nil
}

// GetPGConnectionString ...
func (c *Config) GetPGConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.viper.GetString("postgres.host"),
		c.viper.GetString("postgres.port"),
		c.viper.GetString("postgres.dbname"),
		c.viper.GetString("postgres.user"),
		c.viper.GetString("postgres.password"),
		c.viper.GetString("postgres.sslmode"),
	)
}

// GetRedisConfig ...
func (c *Config) GetRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       0,
	}
}

// GetDockerConfig ...
func (c *Config) GetDockerConfig() map[string]string {
	dockerConfig := map[string]string{
		"host":    c.viper.GetString("docker.host"),
		"version": c.viper.GetString("docker.version"),
		"repo":    c.viper.GetString("docker.repo"),
	}
	return dockerConfig
}

// GetAPIPorts ...
func (c *Config) GetAPIPorts() int {
	return c.viper.GetInt("api.portnum")
}

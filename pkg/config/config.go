package config

import (
	"fmt"

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
	c := Config{viper.New(), env}
	c.viper.SetConfigType("yaml")
	c.viper.SetConfigName("conf")
	c.viper.AddConfigPath(configDir)

	err := c.viper.ReadInConfig()
	return &c, err
}

// GetPGConnectionString ...
func (c *Config) GetPGConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.viper.GetString(c.env+".postgres.host"),
		c.viper.GetString(c.env+".postgres.port"),
		c.viper.GetString(c.env+".postgres.dbname"),
		c.viper.GetString(c.env+".postgres.user"),
		c.viper.GetString(c.env+".postgres.password"),
		c.viper.GetString(c.env+".postgres.sslmode"),
	)
}

// GetRedisConfig ...
func (c *Config) GetRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString(c.env+".redis.host") + ":" + c.viper.GetString(c.env+".redis.port"),
		Password: c.viper.GetString(c.env + ".redis.password"),
		DB:       0,
	}
}

// GetDockerConfig ...
func (c *Config) GetDockerConfig() map[string]string {
	dockerConfig := map[string]string{
		"host":    c.viper.GetString(c.env + ".docker.host"),
		"version": c.viper.GetString(c.env + ".docker.version"),
		"repo":    c.viper.GetString(c.env + ".docker.repo"),
	}
	return dockerConfig
}

// GetAPIPorts ...
func (c *Config) GetAPIPorts() int {
	return c.viper.GetInt(c.env + ".api.portnum")
}

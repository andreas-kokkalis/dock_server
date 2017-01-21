package conf

import (
	"log"

	"github.com/spf13/viper"
)

var v *viper.Viper

// InitConf initializes the viper configuration manager
func InitConf(path string) (err error) {
	v = viper.GetViper()
	v.SetConfigType("yaml")
	v.SetConfigName("conf")
	v.AddConfigPath(path)

	err = v.ReadInConfig()
	if err != nil {
		log.Printf("[InitConf]: could read configuration file.\n")
	}
	return err
}

// GetVal returns a string conf
func GetVal(key string) string {
	return v.GetString(key)
}

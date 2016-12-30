package conf

import (
	"log"

	"github.com/spf13/viper"
)

var v *viper.Viper

// Init the config
func Init() {
	v = viper.GetViper()
	v.SetConfigType("yaml")
	v.SetConfigName("conf")
	v.AddConfigPath("./conf")
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

// GetVal returns a string conf
func GetVal(key string) string {
	return v.GetString(key)
}

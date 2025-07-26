package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Configuration() string {
	viper.SetConfigFile("healthfers/config/config.env")
	viper.ReadInConfig()
	user := viper.Get("user").(string)
	host := viper.Get("host").(string)
	password := viper.Get("password").(string)
	dbname := viper.Get("dbname").(string)
	port := viper.Get("port").(string)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)
	return dsn
}

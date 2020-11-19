package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	port := viper.GetString("port")

	host := viper.GetStringMap("db")

	log.Println(port, host)
}

package main

import (
	"github.com/ggpd/brackets/env"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	router := gin.Default()
	e := env.New()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		e.Log.Fatalf("Error reading config file: %s \n", err)
	}

	url := viper.GetString("app.url")

	user := viper.GetString("database.username")
	pass := viper.GetString("database.password")
	database := viper.GetString("database.name")
	host := viper.GetString("database.host")

	e.ConnectDb(env.DbString(database, host, user, pass))

	e.Log.Printf("Server starting...")
	e.Log.Fatal(autotls.Run(router, url))
}

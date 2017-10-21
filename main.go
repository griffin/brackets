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

	//url := viper.GetString("app.url")

	user := viper.GetString("database.username")
	pass := viper.GetString("database.password")
	database := viper.GetString("database.name")
	host := viper.GetString("database.host")

	e.ConnectDb(env.DbString(database, host, user, pass))

	router.GET("/", nil)

	router.GET("/tournament/:selector", nil)
	router.PUT("/tournament/:selector", nil)
	router.DELETE("/tournament/:selector", nil)
	router.POST("/tournament", nil)

	router.GET("/team/:selector", nil)
	router.PUT("/team/:selector", nil)
	router.DELETE("/team/:selector", nil)
	router.POST("/team", nil)

	router.GET("/user/:selector", nil)
	router.PUT("/user/:selector", nil)
	router.DELETE("/user/:selector", nil)
	router.POST("/user", nil)

	router.GET("/post/:selector", nil)
	router.PUT("/post/:selector", nil)
	router.DELETE("/post/:selector", nil)
	router.POST("/post", nil)

	router.POST("/session", nil)

	e.Log.Printf("Server starting...")
	e.Log.Fatal(autotls.Run(router, "127.0.0.1:4443"))
}

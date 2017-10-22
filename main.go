package main

import (
	"github.com/ggpd/brackets/env"
	"github.com/ggpd/brackets/routes"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"

	"fmt"
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
	port := viper.GetInt("app.port")

	userSql := viper.GetString("sql.username")
	passSql := viper.GetString("sql.password")
	databaseSql := viper.GetString("sql.database")
	hostSql := viper.GetString("sql.host")
	portSql := viper.GetInt("sql.port")

	passRedis := viper.GetString("redis.password")
	hostRedis := viper.GetString("redis.host")


	sqlOptions := env.SqlOptions {
		User: userSql,
		Password: passSql,
		Host: hostSql,
		Port: portSql,
		Database: databaseSql,
	}


	redisOptions := redis.Options{
		Addr: hostRedis,
		Password: passRedis,
		DB: 0,
	}


	e.ConnectDb(sqlOptions, redisOptions)

	router.GET("/", nil)

	router.GET("/login", e.getLoginRoute)
	router.POST("/login", nil)
	router.POST("/logout", nil)

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
	e.Log.Fatal(autotls.Run(router, fmt.Sprintf("%v:%v", url, port)))
}

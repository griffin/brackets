package main

import (
	"fmt"
	"github.com/ggpd/brackets/env"
	"github.com/ggpd/brackets/routes"
	//"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	router := gin.Default()
	e := routes.CastEnv(env.New())

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		e.Log.Fatalf("Error reading config file: %s \n", err)
	}

	url := viper.GetString("app.url")
	port := viper.GetInt("app.port")

	userSQL := viper.GetString("sql.username")
	passSQL := viper.GetString("sql.password")
	databaseSQL := viper.GetString("sql.database")
	hostSQL := viper.GetString("sql.host")
	portSQL := viper.GetInt("sql.port")

	passRedis := viper.GetString("redis.password")
	hostRedis := viper.GetString("redis.host")
	portRedis := viper.GetInt("redis.port")

	sqlOptions := env.SQLOptions{
		User:     userSQL,
		Password: passSQL,
		Host:     hostSQL,
		Port:     portSQL,
		Database: databaseSQL,
	}

	redisOptions := env.RedisOptions{
		Host:     hostRedis,
		Port:     portRedis,
		Password: passRedis,
	}

	e.ConnectDb(sqlOptions, redisOptions)

	_, err = e.Template.ParseFiles("public/index.html",
		"public/notfound.html",
		"public/user/user_index.html",
		"public/team/team_index",
		"public/tournament/tournament_index.html")

	if err != nil {
		e.Log.Fatal(err)
	}

	router.GET("/", e.GetHomeRoute)

	router.GET("/login", e.GetLoginRoute)
	router.POST("/login", nil)
	router.POST("/logout", nil)
	router.GET("/register", nil)

	router.GET("/tournament/:selector", nil)
	router.PUT("/tournament/:selector", nil)
	router.DELETE("/tournament/:selector", nil)
	router.POST("/tournament", nil)

	router.GET("/team/:selector", nil)
	router.PUT("/team/:selector", nil)
	router.DELETE("/team/:selector", nil)
	router.POST("/team", nil)

	router.GET("/user/:selector", e.GetUser)
	router.PUT("/user/:selector", nil)
	router.DELETE("/user/:selector", nil)
	router.POST("/user", nil)

	e.Log.Printf("Server starting...")
	//e.Log.Fatal(autotls.Run(router, fmt.Sprintf("%v:%v", url, port)))
	router.Run(fmt.Sprintf("%v:%v", url, port))

}

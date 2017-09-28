package main

import (
	"github.com/ggpd/teambuilder/env"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	env := env.New()
	env.ConnectDb("")

	env.GetPlayers(1)

	env.Logger.Printf("Server starting...")
	env.Logger.Fatal(autotls.Run(router, "team.oakley.ninja"))
}

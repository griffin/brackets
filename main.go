package main

import (
	"log"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	log.Printf("Server starting...")
	log.Fatal(autotls.Run(router, "team.oakley.ninja"))
}

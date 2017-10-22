package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ggpd/brackets/env"
)

type Env env.Env

func (e *Env) GetLoginRoute(c *gin.Context){
	e.Db.CheckSession("")
}
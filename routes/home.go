package routes

import (
	"github.com/ggpd/brackets/env"
	"github.com/gin-gonic/gin"
)

type Env struct { *env.Env }

func CastEnv(e *env.Env) *Env {
	return &Env{e}
}


func (e *Env) GetHomeRoute(c *gin.Context){

}
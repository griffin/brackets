package routes

import (
	"github.com/gin-gonic/gin"
)

func (e *Env) GetLoginRoute(c *gin.Context){
	e.Db.CheckSession("")
}
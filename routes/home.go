package routes

import (
	"net/http"

	"github.com/ggpd/brackets/env"
	"github.com/gin-gonic/gin"
)

type Env struct{ *env.Env }

func CastEnv(e *env.Env) *Env {
	return &Env{e}
}

func (e *Env) GetHomeRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")
	if err != nil {
		c.HTML(http.StatusOK, "home.html", nil)
		return
	}

	usr, err := e.Db.CheckSession(token)
	if err != nil {
		c.HTML(http.StatusOK, "home.html", nil)
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"user": usr,
	})

}

func NotFoundRoute(c *gin.Context) {
	c.HTML(http.StatusNotFound, "notfound.html", nil)
}

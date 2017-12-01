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
		c.SetCookie("user_session", "del", -1, "/", "", false, false)
		c.HTML(http.StatusOK, "home.html", nil)
		return
	}


	teams, err := e.Db.GetTeams(*usr)
	if err != nil {
		e.Log.Println(err)
	}

	games, err := e.Db.GetUpcomingGames(teams)
	if err != nil {
		e.Log.Println(err)
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"login": usr,
		"teams": teams,
		"games": games,
	})

}

func NotFoundRoute(c *gin.Context) {
	c.HTML(http.StatusNotFound, "notfound.html", nil)
}

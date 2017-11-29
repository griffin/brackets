package routes

import (
	"net/http"

	"github.com/ggpd/brackets/env"
	"github.com/gin-gonic/gin"
)

func (e *Env) GetTeamRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")

	var login *env.User

	if err == nil {
		login, err = e.Db.CheckSession(token)
	}

	team, err := e.Db.GetTeam(c.Param("selector"), true)
	if err != nil {
		e.Log.Println(err)
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}

	var teams []env.Team
	teams = append(teams, *team)

	games, err := e.Db.GetUpcomingGames(teams)

	c.HTML(http.StatusOK, "team_index.html", gin.H{
		"login": login,
		"team":  team,
		"games": games,
	})
}

func (e *Env) PostTeamRoute(c *gin.Context) {

}

package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *Env) GetTeamRoute(c *gin.Context) {
	team, err := e.Db.GetTeam(c.Param("selector"), true)
	if err != nil {
		e.Log.Println(err)
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}

	c.HTML(http.StatusOK, "team_index.html", gin.H{
		"team":  team,
		"games": nil,
	})
}

func (e *Env) PostTeamRoute(c *gin.Context) {

}

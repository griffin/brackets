package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *Env) GetTournamentRoute(c *gin.Context) {
	tour, err := e.Db.GetTournament(c.Param("selector"), true)
	if err != nil {
		e.Log.Println(err)
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}

	c.HTML(http.StatusOK, "tournament_index.html", tour)
}

func (e *Env) PostTournamentRoute(c *gin.Context) {

}

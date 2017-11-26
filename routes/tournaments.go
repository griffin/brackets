package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (e *Env) GetTournamentsRoute(c *gin.Context) {

	pageStr := c.DefaultQuery("page", "0")
	resultsStr := c.DefaultQuery("results", "30")

	page, err1 := strconv.Atoi(pageStr)
	results, err2 := strconv.Atoi(resultsStr)

	if err1 != nil || err2 != nil {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
	}

	tour, left, err := e.Db.GetTournaments(results, page)
	if err != nil {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
	}

	c.HTML(http.StatusOK, "tournaments.html", gin.H{
		"tournaments":    tour,
		"results":        results,
		"pageNumber":     page,
		"nextPageNumber": page + 1,
		"prevPageNumber": page - 1,
		"next":           left > results,
	})
}

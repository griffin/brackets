package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const amountResults = 30

func (e *Env) GetUserRoute(c *gin.Context) {

	usr, err := e.Db.GetUser(c.Param("selector"))
	if err != nil {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_index.html", usr)
}

func (e *Env) GetUsersRoute(c *gin.Context) {

	pageStr := c.DefaultQuery("page", "0")
	resultsStr := c.DefaultQuery("results", "30")

	page, err1 := strconv.Atoi(pageStr)
	results, err2 := strconv.Atoi(resultsStr)

	if err1 != nil || err2 != nil {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
	}

	usr, left, err := e.Db.GetUsers(results, page)
	if err != nil {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
	}

	e.Log.Println(left)

	c.HTML(http.StatusOK, "users.html", gin.H{
		"users":          usr,
		"results":        results,
		"pageNumber":     page,
		"nextPageNumber": page + 1,
		"prevPageNumber": page - 1,
		"next":           left > results,
	})
}

func (e *Env) GetSettingsRoute(c *gin.Context) {

}

func (e *Env) PostSettingsRoute(c *gin.Context) {

}

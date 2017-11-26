package routes

import (
	"net/http"
	"strconv"

	"github.com/ggpd/brackets/env"

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
	token, err := c.Cookie("user_session")
	if err != nil {
		c.Redirect(http.StatusFound, "/")
	}

	usr, err := e.Db.CheckSession(token)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
	}

	c.HTML(http.StatusOK, "user_edit.html", gin.H{
		"user":                usr,
		genderSel(usr.Gender): usr.Gender.String(),
	})
}

func genderSel(g env.Gender) string {
	switch g {
	case env.Male:
		return "male"
	case env.Female:
		return "female"
	case env.Other:
		return "other"
	case env.PreferNotToSay:
		return "pnto"
	}
	return ""
}

func (e *Env) PostSettingsRoute(c *gin.Context) {

}

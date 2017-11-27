package routes

import (
	"net/http"
	"strconv"

	"github.com/ggpd/brackets/env"

	"github.com/gin-gonic/gin"
)


func (e *Env) GetUserRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")

	var login *env.User

	if err == nil {
		login, err = e.Db.CheckSession(token)
	}


	usr, err := e.Db.GetUser(c.Param("selector"))
	if err != nil {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_index.html", gin.H{
		"login": login,
		"user": usr,
	})
}

func (e *Env) GetUsersRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")
	
		var login *env.User
	
		if err == nil {
			login, err = e.Db.CheckSession(token)
		}

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
		"login": login,
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
	token, err := c.Cookie("user_session")
	if err != nil {
		c.Redirect(http.StatusFound, "/")
	}

	usr, err := e.Db.CheckSession(token)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
	}

	first, e1 := c.GetPostForm("first_name")
	last, e2 := c.GetPostForm("last_name")
	email, e3 := c.GetPostForm("email")
	dobSt, e4 := c.GetPostForm("dob")
	genderSt, e5 := c.GetPostForm("gender")

	if !validField(first, e1) ||
		!validField(last, e2) ||
		!validField(email, e3) ||
		!validField(dobSt, e4) ||
		!validField(genderSt, e5) {

		c.HTML(http.StatusOK, "user_edit.html", gin.H{
			"message": "you must complete all forms",
			"type":    "danger",
		})
		return
	}

	err = e.Db.UpdateUser(*usr)
	if err != nil {
		c.HTML(http.StatusOK, "user_edit.html", gin.H{
			"message": "user could not be updated",
			"type":    "danger",
		})
	}

	c.HTML(http.StatusOK, "user_edit.html", gin.H{
		"user": usr,
	})

}

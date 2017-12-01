package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ggpd/brackets/env"

	"github.com/gin-gonic/gin"
)

func (e *Env) GetUserRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")

	var login *env.User

	if err == nil {
		login, err = e.Db.CheckSession(token)
		if err != nil {
			c.SetCookie("user_session", "del", -1, "/", "", false, false)
		}
	}

	usr, err := e.Db.GetUser(c.Param("selector"))
	if err != nil {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_index.html", gin.H{
		"login": login,
		"user":  usr,
	})
}

func (e *Env) GetUsersRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")

	var login *env.User

	if err == nil {
		login, err = e.Db.CheckSession(token)
		if err != nil {
			c.SetCookie("user_session", "del", -1, "/", "", false, false)
		}
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
		"login":          login,
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
		return
	}

	usr, err := e.Db.CheckSession(token)
	if err != nil {

		c.SetCookie("user_session", "del", -1, "/", "", false, false)
		c.Redirect(http.StatusFound, "/")
		return
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
		c.SetCookie("user_session", "del", -1, "/", "", false, false)
		c.Redirect(http.StatusFound, "/")
		return
	}

	first, e1 := c.GetPostForm("first")
	last, e2 := c.GetPostForm("last")
	email, e3 := c.GetPostForm("email")
	dobSt, e4 := c.GetPostForm("dob")
	genderSt, e5 := c.GetPostForm("gender")

	if !validField(first, e1) ||
		!validField(last, e2) ||
		!validField(email, e3) ||
		!validField(dobSt, e4) ||
		!validField(genderSt, e5) {

		c.HTML(http.StatusOK, "user_edit.html", gin.H{
			"user":                usr,
			genderSel(usr.Gender): usr.Gender.String(),
			"message":             "you must complete all forms",
			"type":                "danger",
		})
		return
	}

	usr.FirstName = first
	usr.LastName = last
	usr.Email = email
	usr.Gender = env.ToGender(genderSt)

	dob, err := time.Parse("2006-01-02", dobSt)
	if err != nil {
		c.HTML(http.StatusOK, "user_register.html", gin.H{
			"user":                usr,
			genderSel(usr.Gender): usr.Gender.String(),
			"message":             "date in incorrect format",
			"type":                "danger",
		})
		return
	}
	usr.DateOfBirth = dob

	err = e.Db.UpdateUser(*usr)
	e.Log.Println(err)
	if err != nil {
		c.HTML(http.StatusOK, "user_edit.html", gin.H{
			"user":                usr,
			genderSel(usr.Gender): usr.Gender.String(),
			"message":             "user could not be updated",
			"type":                "danger",
		})
	}

	c.HTML(http.StatusOK, "user_edit.html", gin.H{
		"user":                usr,
		genderSel(usr.Gender): usr.Gender.String(),
	})

}

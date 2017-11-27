package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/ggpd/brackets/env"
	"github.com/gin-gonic/gin"
)

const (
	loginFailed = "Incorrect username or password."
	noUsrOrPsw  = "You need to provide a username and password."
)

func (e *Env) GetLoginRoute(c *gin.Context) {
	_, err := c.Cookie("user_session")
	if err == nil {
		c.Redirect(http.StatusFound, "/")
	}

	c.HTML(http.StatusOK, "user_login.html", nil)
}

func (e *Env) PostLoginRoute(c *gin.Context) {
	email, e1 := c.GetPostForm("email")
	password, e2 := c.GetPostForm("passwd")

	if !e1 || !e2 || len(email) == 0 || len(password) == 0 {
		c.HTML(http.StatusOK, "user_login.html", gin.H{
			"message": noUsrOrPsw,
			"type":    "danger",
		})
		return
	}

	_, token, err := e.Db.CreateSession(email, password)
	if err != nil {
		c.HTML(http.StatusOK, "user_login.html", gin.H{
			"message": loginFailed,
			"type":    "danger",
		})
		return
	}

	e.Log.Println(token)
	c.SetCookie("user_session", token, 120, "/", "", false, false)
	c.Redirect(http.StatusFound, "/")
}

func (e *Env) GetLogoutRoute(c *gin.Context) {
	c.SetCookie("user_session", "del", -1, "/", "", false, false)
	c.Redirect(http.StatusFound, "/")
}

func (e *Env) GetRegisterRoute(c *gin.Context) {
	_, err := c.Cookie("user_session")
	if err == nil {
		c.Redirect(http.StatusFound, "/")
	}

	c.HTML(http.StatusOK, "user_register.html", nil)
}

func (e *Env) PostRegisterRoute(c *gin.Context) {
	first, e1 := c.GetPostForm("first_name")
	last, e2 := c.GetPostForm("last_name")
	email, e3 := c.GetPostForm("email")
	password, e4 := c.GetPostForm("passwd1")
	passwordConf, e5 := c.GetPostForm("passwd2")
	dobSt, e6 := c.GetPostForm("dob")
	genderSt, e7 := c.GetPostForm("gender")

	if !validField(first, e1) ||
		!validField(last, e2) ||
		!validField(email, e3) ||
		!validField(password, e4) ||
		!validField(passwordConf, e5) ||
		!validField(dobSt, e6) ||
		!validField(genderSt, e7) {

		c.HTML(http.StatusOK, "user_register.html", gin.H{
			"message": "you must complete all forms",
			"type":    "danger",
		})
		return
	}

	g := env.ToGender(genderSt)

	dob, err := time.Parse("2006-01-02", dobSt)
	if err != nil {
		c.HTML(http.StatusOK, "user_register.html", gin.H{
			"message": "date in incorrect format",
			"type":    "danger",
		})
		return
	}

	if strings.Compare(password, passwordConf) != 0 {
		c.HTML(http.StatusOK, "user_register.html", gin.H{
			"message": "Passwords do not match.",
			"type":    "danger",
		})
		return
	}

	usr := env.User{
		FirstName:   first,
		LastName:    last,
		Email:       email,
		Gender:      g,
		DateOfBirth: dob,
	}

	_, err = e.Db.CreateUser(usr, password)
	if err != nil {
		c.HTML(http.StatusOK, "user_register.html", gin.H{
			"message": "error creating user",
			"type":    "danger",
		})
		return
	}

	_, token, err := e.Db.CreateSession(usr.Email, password)
	if err != nil {
		c.HTML(http.StatusOK, "user_register.html", gin.H{
			"message": "error creating session",
			"type":    "danger",
		})
		return
	}
	
	c.SetCookie("user_session", token, 120, "/", "", false, false)
	c.Redirect(http.StatusFound, "/")
}

func validField(field string, rec bool) bool {
	return rec && len(field) > 0
}

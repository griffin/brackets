package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	loginFailed = "Incorrect username or password."
	noUsrOrPsw  = "You need to provide a username and password."
)

func (e *Env) GetLoginRoute(c *gin.Context) {
	_, err := c.Cookie("session")
	if err == nil {
		c.Redirect(302, "/")
	}

	c.HTML(http.StatusOK, "user_login.html", gin.H{
		"error": "",
	})

}

func (e *Env) PostLoginRoute(c *gin.Context) {
	email, e1 := c.GetPostForm("email")
	password, e2 := c.GetPostForm("passwd")

	if !e1 || !e2 || len(email) == 0 || len(password) == 0 {
		c.HTML(http.StatusOK, "user_login.html", gin.H{
			"message": noUsrOrPsw,
			"type": "danger",
		})
		return
	}

	_, token, err := e.Db.CreateSession(email, password)
	if err != nil {
		c.HTML(http.StatusOK, "user_login.html", gin.H{
			"message": loginFailed,
			"type": "danger",
		})
		return
	}

	c.SetCookie("session", token, -1, "", "", false, false)
	c.Redirect(302, "/")
}

func (e *Env) PostLogoutRoute(c *gin.Context) {

}

func (e *Env) GetRegisterRoute(c *gin.Context) {

}

func (e *Env) PostRegisterRoute(c *gin.Context) {

}

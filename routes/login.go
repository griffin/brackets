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
	e.Log.Println("here0")
	if !e1 || !e2 {
		c.HTML(http.StatusOK, "user_login.html", gin.H{
			"error": noUsrOrPsw,
		})
		return
	}
	e.Log.Println("here1")
	_, token, err := e.Db.CreateSession(email, password)
	if err != nil {
		c.HTML(http.StatusOK, "user_login.html", gin.H{
			"error": loginFailed,
		})
		return
	}

	e.Log.Println("here2")
	c.SetCookie("session", token, -1, "", "", false, false)
	c.Redirect(302, "/")
}

func (e *Env) PostLogoutRoute(c *gin.Context) {

}

func (e *Env) GetRegisterRoute(c *gin.Context) {

}

func (e *Env) PostRegisterRoute(c *gin.Context) {

}

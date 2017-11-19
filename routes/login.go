package routes

import (
	"github.com/gin-gonic/gin"
)

func (e *Env) GetLoginRoute(c *gin.Context) {
	_, err := c.GetCookie("session")
	if err == nil {
		c.Redirect(302, "/")
	}

}

func (e *Env) PostLoginRoute(c *gin.Context) {
	email, e := c.GetPostForm("email")
	password, e := c.GetPostForm("passwd")

	_, token, err := e.Db.CreateSession(email, password)
	if err != nil {

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

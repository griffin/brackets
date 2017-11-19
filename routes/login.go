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
	email, e1 := c.GetPostForm("email")
	password, e2 := c.GetPostForm("passwd")
	if e1 || e2 {

	}

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

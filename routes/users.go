package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *Env) GetUserRoute(c *gin.Context) {

	t := e.Template.Lookup("user_index.html")

	usr, err := e.Db.GetUser(c.Param("selector"))
	if err != nil {
		t = e.Template.Lookup("notfound.html")
		c.Status(http.StatusNotFound)
		t.Execute(c.Writer, nil)
		return
	}

	c.Status(http.StatusOK)
	t.Execute(c.Writer, usr)
}

func (e *Env) PostUserRoute(c *gin.Context) {

}

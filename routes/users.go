package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *Env) GetUserRoute(c *gin.Context) {

	usr, err := e.Db.GetUser(c.Param("selector"))
	if err != nil {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_index.html", usr)
}

func (e *Env) GetUsersRoute(c *gin.Context) {

}

func (e *Env) GetSettingsRoute(c *gin.Context) {

}

func (e *Env) PostSettingsRoute(c *gin.Context) {

}

package routes

import (
	"github.com/ggpd/brackets/env"
	"github.com/gin-gonic/gin"
)

const (
	bootStrap = ""
)

type Env struct { *env.Env }

func CastEnv(e *env.Env) *Env {
	return &Env{e}
}


func (e *Env) GetHomeRoute(c *gin.Context){

}

func PushJS() gin.HandlerFunc {
	return func(c *gin.Context) {

		if pusher, ok := c.Writer.(http.Pusher); ok {
            options := &http.PushOptions{
                Header: http.Header{
                    "Accept-Encoding": r.Header["Accept-Encoding"],
                },
			}
			
            if err := pusher.Push(bootStrap, options); err != nil {
                log.Printf("Failed to push: %v", err)
			}
		}
		
	}
}
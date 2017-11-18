package routes

import (
	"html/template"
	"io/ioutil"
)

func (e *Env) GetUser(c *gin.Context) {

	src, _ := ioutil.ReadFile("../templates/user/userprofile.html")
	t, _ := template.New("userprofile").Parse(string(src))

	t.Execute(c.Writer, nil)

}

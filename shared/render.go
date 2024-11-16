package shared

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, code int, templateName string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	globalData := c.MustGet("templateData").(gin.H)
	for k, v := range globalData {
		data[k] = v
		fmt.Println(k, v)
	}
	c.HTML(code, templateName, data)
}

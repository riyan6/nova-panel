// internal/webserver/router.go
package webserver

import (
	"github.com/gin-gonic/gin"
	"nova-panel/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/agents", v1.GetAgentList)
		// 你可以继续添加 api.POST("/command") 等
	}

	return r
}

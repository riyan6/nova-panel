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
	}

	// 初始化 WebSocket 服务
	initWebSocket(r)

	return r
}

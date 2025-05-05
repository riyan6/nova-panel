// api/v1/agent.go
package v1

import (
	"net/http"
	"nova-panel/internal/store"

	"github.com/gin-gonic/gin"
)

func GetAgentList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"count":  store.GetAgentCount(),
		"agents": store.GetAllAgents(),
	})
}

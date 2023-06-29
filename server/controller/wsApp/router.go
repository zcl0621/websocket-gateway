package wsApp

import (
	"github.com/gin-gonic/gin"
)

func InitWsAppRouter(api *gin.RouterGroup) {
	higoApp := api.Group("app")
	{
		higoApp.GET("ws", wsApp)
	}
	api.GET("group/msg/index", groupIndexMsg)
	api.POST("group/msg/index", groupIndexMsg)
	api.GET("user/msg/index", userIndexMsg)
	api.POST("user/msg/index", userIndexMsg)
}

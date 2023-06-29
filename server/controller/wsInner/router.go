package wsInner

import "github.com/gin-gonic/gin"

func InitWsInnerRouter(api *gin.RouterGroup) {
	api.POST("group-publish", groupPublish)
	api.POST("user-publish", userPublish)
	api.POST("client-publish", clientPublish)
}

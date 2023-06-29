package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"net/http"
	"time"
	"ws-gateway/config"
	"ws-gateway/server/controller/wsApp"
	"ws-gateway/server/controller/wsInner"
	"ws-gateway/wsManager"

	ginSwagger "github.com/swaggo/gin-swagger"
	_ "ws-gateway/docs"
	"ws-gateway/exception"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(exception.ExceptionHandler)
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowWebSockets:  true,
		AllowAllOrigins:  true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	wsManager.DoInit()
	go wsManager.Manager.Start()
	go wsManager.WriteMessage()
	go wsManager.GroupChannelsManager.Start()

	publicGroup := r.Group("/api/public/ws-gateway-service")
	wsApp.InitWsAppRouter(publicGroup)

	innerGroup := r.Group("/api/inner/ws-gateway-service")
	wsInner.InitWsInnerRouter(innerGroup)

	if config.RunMode == "debug" || config.RunMode == "dev" {
		r.GET("/swagger/ws-gateway-service/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "")
	})
	return r
}

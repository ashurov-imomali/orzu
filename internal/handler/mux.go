package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func InitRoutes(h Handler, token string) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.Use(gin.Recovery(), h.Logger(), h.ServiceAuth(token))

	e.GET("/ping", h.ping)

	{
		user := e.Group("/user")
		user.GET("/check-pan", h.checkUserPan) // 5
		user.GET("/:inn", h.getUserByInn)      // 1
	}

	{
		otp := e.Group("/otp")
		otp.POST("/send", h.sendOtp)           // 2
		otp.POST("/confirm-otp", h.confirmOtp) // 3
	}

	{
		business := e.Group("/business")
		business.Use(h.Auth())
		business.GET("/services", h.getServices)
		business.GET("/conditions", h.getConditions) //4
		business.POST("/transh", h.createTransh)     //6
	}

	return e
}

package http

import (
	_ "github.com/martinusiron/loan-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(h *Handler) *gin.Engine {
	r := gin.Default()
	NewHandler(r, h.UC)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

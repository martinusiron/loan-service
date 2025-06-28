package http

import (
	_ "github.com/martinusiron/loan-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/martinusiron/loan-service/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(uc *usecase.LoanUsecase) *gin.Engine {
	r := gin.Default()
	NewHandler(r, uc)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

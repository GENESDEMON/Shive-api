package routes

import (
	controller "github.com/genesdemon/golang-jwt-project/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/signin", controller.Login())
}

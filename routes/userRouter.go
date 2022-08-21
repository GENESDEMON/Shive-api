package routes

import (
	controller "github.com/genesdemon/golang-jwt-project/controllers"
	"github.com/genesdemon/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.PUT("/users/edituser", controller.EditUser())
}

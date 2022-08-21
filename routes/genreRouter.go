package routes

import (
	"github.com/genesdemon/golang-jwt-project/controllers"
	"github.com/genesdemon/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func GenreRoutes(incomingRoutes gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/admin/creategenre", controllers.CreateGenre())
	incomingRoutes.GET("/admin/:genre_id", controllers.GetGenre())
	incomingRoutes.GET("/admin/getgenres", controllers.GetGenres())
	incomingRoutes.PUT("/admin/editgenre", controllers.EditGenre())
	//incomingRoutes.PUT("/admin/deletegenre", controllers.DeleteGenre())

}

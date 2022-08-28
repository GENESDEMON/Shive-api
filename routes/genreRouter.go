package routes

import (
	"github.com/genesdemon/golang-jwt-project/controllers"
	"github.com/genesdemon/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func GenreRoutes(incomingRoutes gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/genres/creategenre", controllers.CreateGenre())
	incomingRoutes.GET("/genres/:genre_id", controllers.GetGenre())
	incomingRoutes.GET("/genres/getgenres", controllers.GetGenres())
	incomingRoutes.PUT("/genres/editgenre/:genre_id", controllers.EditGenre())
	incomingRoutes.DELETE("/genres/:genre_id", controllers.DeleteAGenre())

}

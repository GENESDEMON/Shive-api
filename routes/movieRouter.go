package routes

import (
	"github.com/genesdemon/golang-jwt-project/controllers"
	"github.com/genesdemon/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func MovieRoutes(incomingRoutes gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/movies/createmovie", controllers.CreateMovie())
	incomingRoutes.GET("/movies/:movie_id", controllers.GetMovie())
	incomingRoutes.GET("/movies/getmovies", controllers.GetMovies())
	incomingRoutes.PUT("/movies/editmovie/:movie_id", controllers.EditMovie())
	incomingRoutes.DELETE("/movies/:movie_id", controllers.DeleteMovie())
	incomingRoutes.GET("/movies/search", controllers.SearchMovieByQuery())
}

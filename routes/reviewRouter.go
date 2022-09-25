package routes

import (
	"github.com/genesdemon/golang-jwt-project/controllers"
	"github.com/genesdemon/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func ReviewRoutes(incomingRoutes gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("reviews/addreview", controllers.AddAReview())
	incomingRoutes.DELETE("/reviews/:_id", controllers.DeleteAReview())
	incomingRoutes.GET("/movies/:movie_id", controllers.ViewAMovieReviews())
	incomingRoutes.GET("/movies/:reviewer_id", controllers.AllUserReviews())
}

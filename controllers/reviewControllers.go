package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/genesdemon/golang-jwt-project/database"
	helper "github.com/genesdemon/golang-jwt-project/helpers"
	"github.com/genesdemon/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewCollection *mongo.Collection = database.OpenCollection(database.Client, "review")

//Add  new review
func AddAReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "USER"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var review models.Reviews
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}
		//use the validator library to validate required fields
		if validationErr := validate.Struct(&review); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newReview := models.Reviews{
			Id:          primitive.NewObjectID(),
			Movie_id:    review.Movie_id,
			Reviewer_id: review.Reviewer_id,
			Review:      review.Review,
		}

		result, err := reviewCollection.InsertOne(ctx, newReview)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Message": "success",
			"Data":    map[string]interface{}{"data": result}})
	}
}

//View the reviews of a specific movie
func ViewAMovieReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchreviews []models.Reviews
		queryParam := c.Query("movie_id")
		if queryParam == "" {
			log.Println("No movie id passed")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := reviewCollection.Find(ctx, bson.M{"movie_id": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchquerydb.All(ctx, &searchreviews)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchreviews)
	}
}

//Delete a review
func DeleteAReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		reviewId := c.Param("_id")
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(reviewId)

		result, err := reviewCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				gin.H{
					" Status":  http.StatusNotFound,
					" Message": "error",
					" Data":    map[string]interface{}{"data": "Review with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			gin.H{
				"Status":  http.StatusOK,
				"Message": "success",
				"Data":    map[string]interface{}{"data": "Your review was successfully deleted!"}},
		)
	}
}

// Allow a user view all their Reviews
func AllUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchreviews []models.Reviews
		queryParam := c.Query("reviewer_id")
		if queryParam == "" {
			log.Println("No reviewer id passed")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := reviewCollection.Find(ctx, bson.M{"reviewer_id": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchquerydb.All(ctx, &searchreviews)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchreviews)
	}
}

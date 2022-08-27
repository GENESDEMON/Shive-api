package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/genesdemon/golang-jwt-project/database"
	helper "github.com/genesdemon/golang-jwt-project/helpers"
	"github.com/genesdemon/golang-jwt-project/models"
	"github.com/genesdemon/golang-jwt-project/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var genreCollection *mongo.Collection = database.OpenCollection(database.Client, "genre")

func CreateGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var genre models.Genre
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, responses.GenreResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()}})
			return
		}

		//Check to see if name exists
		regexMatch := bson.M{"$regex": primitive.Regex{Pattern: *genre.Name, Options: "i"}}
		count, err := genreCollection.CountDocuments(ctx, bson.M{"name": regexMatch})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while checking for the Name"})
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this genre name already exists", "count": count})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&genre); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GenreResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newGenre := models.Genre{
			Id:   primitive.NewObjectID(),
			Name: genre.Name,
		}

		result, err := genreCollection.InsertOne(ctx, newGenre)

		if err != nil {
			c.JSON(http.StatusBadRequest, responses.GenreResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()}})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GenreResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.GenreResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data:    map[string]interface{}{"data": result}})
	}
}

// To get just one genre
func GetGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		genreId := c.Param("genre_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var genre models.Genre
		err := genreCollection.FindOne(ctx, bson.M{"genre_id": genreId}).Decode(&genre)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, genre)
	}
}

// To fetch all genres
func GetGenres() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "genre_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := genreCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching genres "})
		}
		var allgenres []bson.M
		if err = result.All(ctx, &allgenres); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allgenres[0])
	}
}

// Edit genre
func EditGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		genreId := c.Param("genre_id")
		var genre models.Genre
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(genreId)

		//validate the request body
		if err := c.BindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&genre); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": genre.Name}
		filterByID := bson.M{"_id": bson.M{"$eq": objId}}
		result, err := genreCollection.UpdateOne(ctx, filterByID, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}
		//get updated genre details
		var updatedGenre models.Genre
		if result.MatchedCount == 1 {
			err := genreCollection.FindOne(ctx, filterByID).Decode(&updatedGenre)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Status":  http.StatusInternalServerError,
					"Message": "error",
					"Data":    map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    updatedGenre})
		// "Data":    map[string]interface{}{"data": updatedGenre}})
	}
}

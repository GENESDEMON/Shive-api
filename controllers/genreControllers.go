package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/genesdemon/golang-jwt-project/database"
	helper "github.com/genesdemon/golang-jwt-project/helpers"
	"github.com/genesdemon/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var genreCollection *mongo.Collection = database.OpenCollection(database.Client, "genre")

//To create a single genre
func CreateGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var genre models.Genre
		if err := c.BindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(genre)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		count, err := genreCollection.CountDocuments(ctx, bson.M{"name": genre.Name})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the Name"})
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this genre name already exists"})
		}

		genre.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		genre.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		genre.Genre_ID = primitive.NewObjectID()
		genre.Genre_id = genre.Genre_ID.Hex()

		resultInsertionNumber, insertErr := genreCollection.InsertOne(ctx, genre)
		if insertErr != nil {
			msg := fmt.Sprintf("Genre was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"MongoID": resultInsertionNumber,
			"message": string("Request complete"), // cast it to string before showing
		})
	}
}

//To get just one genre
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

//To fetch all genres
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

//To edit just one genre
func EditGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, err)
		}
		var editgenre models.Genre
		if err := c.BindJSON(&editgenre); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "genre.name", Value: editgenre.Name}}}}
		_, err = genreCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something Went Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Updated the Genre")
	}
}

//Delete genre

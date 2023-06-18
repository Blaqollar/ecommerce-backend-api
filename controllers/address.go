package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Blaqollar/ecommerce-backend-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("_id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invaid code"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(500, "Invalid user_id")
		}

		var addresses models.Address

		addresses.AddressID = primitive.NewObjectID().Hex()
		err = c.BindJSON(&addresses)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$addressID"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.JSON(500, "Internal server error")
		}

		var addressInfo bson.M
		err = cursor.All(ctx, &addressInfo)
		if err != nil {
			log.Panic(err)
		}

		var size int64
		for _, addressNo := range addressInfo {
			count := addressNo["count"]
			size = count.Int64()
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.JSON(400, "Not Allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_Id := c.Query("_id")

		if user_Id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "UserId is empty"})
			c.Abort()
		}

		address := make([]models.Address, 0)

		userID, err := primitive.ObjectIDFromHex(user_Id)
		if err != nil {
			c.JSON(404, "Internal error")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.M{"_id": userID}
		update := bson.M{"$set": bson.M{"address": address}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(404, "Address not updated")
		}
		defer cancel()
		ctx.Done()

		c.JSON(200, "Successfully deleted address")
	}
}

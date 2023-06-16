package controllers

import (
	"net/http"
	"time"

	"github.com/Blaqollar/ecommerce-backend-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

func AddAddress() gin.HandlerFunc {

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

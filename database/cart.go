package database

import (
	"context"
	"errors"
	"log"

	"github.com/Blaqollar/ecommerce-backend-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId, userID string) error {
	product, err := prodCollection.Find(ctx, bson.M{"_id": productId})
	if err != nil {
		log.Println(err)
		return errors.New("could not find product")
	}

	var productCart []models.ProductUser

	err = product.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return errors.New("could not decode product")
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("user ID not valid")
	}

	filter := bson.M{"_id": id}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("can't update user")
	}
	return nil
}

func RemoveItem() {

}

func BuyFromCart() {

}

func InstantBuyer() {

}

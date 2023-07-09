package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Blaqollar/ecommerce-backend-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// This function adds products to cart
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

// This function removes items from the collection
func RemoveItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, userID, productID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("userID is not found")
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return errors.New("can't remove item from cart")
	}
	return nil
}

// This function creates an order from cart
func BuyFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("userID is not a valid")
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.OrderId = primitive.NewObjectID().Hex()
	orderCart.OrderedAt = time.Now()
	orderCart.OrderCart = make([]models.ProductUser, 0)
	orderCart.PaymentMethod.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$userCart.price"}}}}}}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	if err != nil {
		log.Panic(err)
	}
	ctx.Done()

	var getUserCart []bson.M

	err = result.All(ctx, &getUserCart)
	if err != nil {
		log.Panic(err)
	}

	var total_price float64
	for _, userItem := range getUserCart {
		price := userItem["total"]
		total_price = price.(float64)
	}
	orderCart.Price = total_price

	filter := bson.M{"_id": id}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.M{"_id": id}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	userCartEmpty := make([]models.ProductUser, 0)
	filter3 := bson.M{"_id": id}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "userCart", Value: userCartEmpty}}}}

	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// This function creates an instant order
func InstantBuy(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("userID is not a valid")
	}

	var productDetails models.ProductUser
	var orderDetails models.Order

	orderDetails.OrderId = primitive.NewObjectID().Hex()
	orderDetails.OrderedAt = time.Now()
	orderDetails.OrderCart = make([]models.ProductUser, 0)
	orderDetails.PaymentMethod.COD = true

	err = prodCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
	}

	orderDetails.Price = productDetails.Price

	filter := bson.M{"_id": id}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderDetails}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter1 := bson.M{"_id": id}
	update1 := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}

	_, err = userCollection.UpdateOne(ctx, filter1, update1)
	if err != nil {
		log.Println(err)
	}
	return nil
}

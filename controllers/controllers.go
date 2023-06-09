package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Blaqollar/ecommerce-backend-api/database"
	"github.com/Blaqollar/ecommerce-backend-api/models"
	"github.com/Blaqollar/ecommerce-backend-api/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")

// This function generates the hash of the password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// This function verifies the password
func VerifyPassword(userPassword, givenPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	if err != nil {
		return false
	}
	return true
}

// This is the Signup function
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := validator.New()
		err = validate.Struct(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phone number is already in use"})
			return
		}

		user.ID = primitive.NewObjectID().Hex()
		password := HashPassword(user.Password)
		user.UserId = user.ID
		user.Password = password
		user.CreatedAt = time.Now().UTC()
		user.UpdatedAt = time.Now().UTC()

		token, refreshtoken, _ := tokens.TokenGenerator(user.Email, user.FirstName, user.LastName, user.UserId)
		user.Token = token
		user.RefreshToken = refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.AddressDetails = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order, 0)

		_, err = UserCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, "successfully created")
	}
}

// This is the login function
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		filter := bson.M{"email": user.Email}
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		var foundUser models.User
		err = UserCollection.FindOne(ctx, filter).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		PasswordIsValid := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "iinvalid credentials"})
			fmt.Println("invalid credentials")
			return
		}

		token, refreshtoken, _ := tokens.TokenGenerator(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.UserId)
		defer cancel()

		tokens.UpdateAllTokens(token, refreshtoken, foundUser.UserId)

		c.JSON(http.StatusFound, foundUser)
	}
}

// This function adds products by admin
func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products models.Product
		err := c.BindJSON(&products)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		products.ProductID = primitive.NewObjectID().Hex()
		_, err = ProductCollection.InsertOne(ctx, products)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Item not inserted"})
		}

		c.JSON(http.StatusOK, "successfully added")
	}
}

// This function searches the database for products
func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, "something went wrong")
		}

		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.JSON(400, "invalid")
		}
		defer cancel()
		c.JSON(200, productList)
	}
}

// This function searches for a product by name
func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var queryProducts []models.Product

		productQueryID := c.Query("product_Name")

		if productQueryID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"product_Name": "product name is required"})
			c.Abort()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.M{"productName": productQueryID})

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, "something went wrong")
		}

		err = cursor.All(ctx, &queryProducts)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.JSON(400, "invalid")
		}

		defer cancel()
		c.JSON(200, queryProducts)
	}
}

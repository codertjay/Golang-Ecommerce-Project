package controllers

import (
	"Golang-Ecommerce-Project/database"
	"Golang-Ecommerce-Project/models"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection *mongo.Collection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(400, errors.New("product id is empty"))
			return
		}
		userQueryId := c.Query("user_id")
		if userQueryId == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(400, errors.New("user id is empty"))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": err})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productId, userQueryId)
		if err != nil {
			c.IndentedJSON(500, err.Error())
			return
		}
		c.JSON(200, "Successfully added to cart")
	}
}
func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(400, errors.New("product id is empty"))
			return
		}
		userQueryId := c.Query("user_id")
		if userQueryId == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(400, errors.New("user id is empty"))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": err})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productId, userQueryId)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, err)
			return
		}
		c.IndentedJSON(200, "Successfully remove item from cart")
	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			log.Println("user id is empty")
			c.JSON(404, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}
		_user_id, _ := primitive.ObjectIDFromHex(user_id)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User
		err := userCollection.FindOne(
			ctx,
			bson.D{primitive.E{
				Key: "_id", Value: _user_id},
			}).Decode(&filledCart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(404, "not found")
			return
		}
		filter_match := bson.D{{
			Key: "$match", Value: bson.D{
				primitive.E{
					Key: "_id", Value: _user_id},
			}}}
		unwind := bson.D{
			{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$user_cart"},
			}}}
		grouping := bson.D{
			{Key: "$group", Value: bson.D{primitive.E{
				Key:   "_id",
				Value: "$_id",
			},
				{Key: "total", Value: bson.D{primitive.E{
					Key: "$sum", Value: "$user_cart.price",
				}}},
			}},
		}
		pointCursor, err := userCollection.Aggregate(
			ctx, mongo.Pipeline{
				filter_match, unwind, grouping,
			})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		err = pointCursor.All(ctx, &listing)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
		}
		log.Println(listing)
		for _, json := range listing {
			c.IndentedJSON(200, gin.H{"total": json["total"], "products": filledCart.UserCart})
		}
		if len(listing) == 0 {
			c.IndentedJSON(200, gin.H{"total": 0, "products": filledCart.UserCart})
		}
		ctx.Done()
	}
}
func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("id")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(400, errors.New("user Id is empty"))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId)
		if err != nil {
			c.IndentedJSON(500, err.Error())
			return
		}
		c.IndentedJSON(200, "Successfully placed the order")
		return
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(400, errors.New("product id is empty"))
			return
		}
		userQueryId := c.Query("user_id")
		if userQueryId == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(400, errors.New("user id is empty"))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": err})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.InstantBuy(ctx, app.prodCollection, app.userCollection, productId, userQueryId)
		if err != nil {
			c.IndentedJSON(500, err)
			return
		}
		c.IndentedJSON(200, "Successfully placed the order")
	}
}

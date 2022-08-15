package controllers

import (
	"Golang-Ecommerce-Project/models"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		user_id := c.Query("user_id")
		if user_id == "" {
			c.Header("Content-Type", "Application/Json")
			c.JSON(404, gin.H{"error": "Invalid Code"})
			c.Abort()
			return
		}
		// the address is the _user_id i have been using i suppose to find
		//one with the _user_id to check if the user exists before filtering match
		_user_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.JSON(500, "Internal server error")
			return
		}
		var addresses models.Address
		addresses.AddressId = primitive.NewObjectID()
		err = c.BindJSON(&addresses)
		if err != nil {
			c.IndentedJSON(406, err.Error())
			return
		}
		matchFilter := bson.D{
			{Key: "$match", Value: bson.D{
				primitive.E{Key: "_id", Value: _user_id},
			}}}
		unwind := bson.D{
			{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$address"},
			}},
		}
		group := bson.D{
			{Key: "$group", Value: bson.D{
				primitive.E{Key: "_id", Value: "$address_id"},
				{Key: "count", Value: bson.D{
					primitive.E{Key: "$sum", Value: 1}}},
			}},
		}
		pointCursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchFilter, unwind, group,
		})
		if err != nil {
			log.Println(err.Error())
			c.IndentedJSON(500, err.Error())
			return
		}
		var addressInfo []bson.M
		err = pointCursor.All(ctx, &addressInfo)
		if err != nil {
			panic(err)
			return
		}
		var size int32
		for _, addressNo := range addressInfo {
			count := addressNo["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: _user_id}}
			update := bson.D{primitive.E{Key: "$push", Value: bson.D{
				primitive.E{Key: "address", Value: addresses},
			}}}
			_, err := userCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				log.Println(err)
				c.IndentedJSON(500, "Unable to Update address")
				return
			}
			c.IndentedJSON(200, " Successfully added Address")
		} else {
			c.IndentedJSON(400, "Not allowed")
		}
		defer cancel()
		ctx.Done()

	}

}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			log.Println("User id is empty")
			c.JSON(404, gin.H{"error": "Invalid search Index"})
			return
		}
		var editAddress models.Address
		err := c.BindJSON(&editAddress)
		if err != nil {
			c.IndentedJSON(400, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		_user_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, gin.H{"error": "Internal server error"})
			return
		}
		filter := bson.D{primitive.E{
			Key:   "_id",
			Value: _user_id,
		}}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				primitive.E{Key: "address.0.house_name", Value: editAddress.House},
				{Key: "address.0.street_name", Value: editAddress.Street},
				{Key: "address.0.city_name", Value: editAddress.City},
				{Key: "address.0.pin_code", Value: editAddress.PinCode},
			}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the home address")
	}

}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			log.Println("User id is empty")
			c.JSON(404, gin.H{"error": "Invalid search Index"})
			return
		}
		var editAddress models.Address
		err := c.BindJSON(&editAddress)
		if err != nil {
			c.IndentedJSON(400, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		_user_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, gin.H{"error": "Internal server error"})
			return
		}
		filter := bson.D{primitive.E{
			Key:   "_id",
			Value: _user_id,
		}}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				primitive.E{Key: "address.1.house_name", Value: editAddress.House},
				{Key: "address.1.street_name", Value: editAddress.Street},
				{Key: "address.1.city_name", Value: editAddress.City},
				{Key: "address.1.pin_code", Value: editAddress.PinCode},
			}}}

		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the work address")

	}

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			log.Println("User id is empty")
			c.JSON(404, gin.H{"error": "Invalid search Index"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		addresses := make([]models.Address, 0)
		_user_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, gin.H{"error": "Internal server error"})
			return
		}
		filter := bson.D{primitive.E{
			Key:   "_id",
			Value: _user_id,
		}}
		// update with an empty value
		update := bson.D{{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "address", Value: addresses},
			}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, gin.H{"error": "Time out"})
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Deleted")

	}

}

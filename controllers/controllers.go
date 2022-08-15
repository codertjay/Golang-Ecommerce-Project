package controllers

import (
	"Golang-Ecommerce-Project/database"
	"Golang-Ecommerce-Project/models"
	"Golang-Ecommerce-Project/tokens"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var validate = validator.New()
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "Users")
var productCollection *mongo.Collection = database.OpenCollection(database.Client, "Products")

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panicln(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword))
	PasswordIsValid := true
	msg := "Password was valid"
	if err != nil {
		PasswordIsValid = false
		msg = "Password was invalid"
		log.Println(err)
		return PasswordIsValid, msg
	}
	return PasswordIsValid, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		if count > 0 {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Email already exist"})
			return
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		if count > 0 {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Phone number already exist"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		user.CreatedAT, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAT, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
		token, refreshToken, err := tokens.TokenGenerator(
			*user.Email, *user.FirstName,
			*user.LastName, *&user.UserId,
		)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.AddressDetails = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order, 0)
		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(500, gin.H{"error": "the user did not get created"})
			return
		}
		defer cancel()
		c.JSON(201, gin.H{"message": "Successfully signed in"})

	}

}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		// user from post request
		var user models.User
		// user from our database if it exists
		var fountUser models.User

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&fountUser)
		defer cancel()
		if err != nil {
			c.JSON(400, gin.H{"error": "login or password incorrect"})
			return
		}
		PasswordIsValid, msg := VerifyPassword(*fountUser.Password, *user.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON(500, gin.H{"error": msg})
			log.Println(msg)
			return
		}
		token, refreshToken, _ := tokens.TokenGenerator(
			*fountUser.Email,
			*fountUser.FirstName,
			*fountUser.LastName,
			fountUser.UserId,
		)
		user.Token = &token
		user.RefreshToken = &refreshToken
		tokens.UpdateAllTokens(token, refreshToken, fountUser.UserId)
		c.JSON(200, fountUser)
	}
}

func ProductViewAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		products.ProductId = primitive.NewObjectID()
		_, anyerr := productCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(400, gin.H{"error": "not inserted"})
			return
		}
		defer cancel()
		c.JSON(200, "Successfully added ")

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var productList []models.Product
		cursor, err := productCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(500, "Something went wrong")
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productList)

	}
}
func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var SearchProducts []models.Product
		queryParam := c.Query("name")

		//	you want to check if its empty
		if queryParam == "" {
			log.Println("Query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(404, gin.H{"error": "Invalid search index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchQueryDB, err := productCollection.Find(ctx, bson.M{
			"product_name": bson.M{"$regex": queryParam},
		})
		if err != nil {
			log.Println(err)
			c.JSON(404, gin.H{"error": "Something went wrong while fetching the data"})
			return
		}
		err = searchQueryDB.All(ctx, &SearchProducts)
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid"})
			return
		}
		defer searchQueryDB.Close(ctx)
		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(500, err)
			return
		}
		defer cancel()
		c.IndentedJSON(200, SearchProducts)

	}
}

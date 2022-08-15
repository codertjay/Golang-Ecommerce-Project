package main

import (
	"Golang-Ecommerce-Project/controllers"
	"Golang-Ecommerce-Project/database"
	"Golang-Ecommerce-Project/middleware"
	"Golang-Ecommerce-Project/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(
		database.OpenCollection(database.Client, "Products"),
		database.OpenCollection(database.Client, "Users"),
	)

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", app.GetItemFromCart())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())

	log.Fatalln(router.Run(":" + port))

}

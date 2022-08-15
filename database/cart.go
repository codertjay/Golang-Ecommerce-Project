package database

import (
	"Golang-Ecommerce-Project/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

var (
	ErrCantFindProduct    = errors.New("cant find the product")
	ErrCantDecodeProducts = errors.New("cant decode the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart(
	ctx context.Context,
	prodCollection *mongo.Collection,
	userCollection *mongo.Collection,
	productId primitive.ObjectID,
	userId string) error {

	searchFromDb, err := prodCollection.Find(
		ctx, bson.M{"_id": productId},
	)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	// UserCart is the productUser struct on the models the
	// UserCart is a value under the User struct

	var productCart []models.ProductUser
	err = searchFromDb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}
	// getting the _id using the string UserId under the user struct
	_user_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{
		primitive.E{Key: "_id", Value: _user_id},
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			primitive.E{Key: "user_cart", Value: bson.D{
				{Key: "$each", Value: productCart},
			}},
		}},
	}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil

}

func RemoveCartItem(
	ctx context.Context,
	prodCollection *mongo.Collection,
	userCollection *mongo.Collection,
	productId primitive.ObjectID,
	userId string) error {
	_user_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{
		primitive.E{Key: "_id", Value: _user_id},
	}
	update := bson.M{
		"$pull": bson.M{"user_cart": bson.M{"_id": productId}},
	}
	// todo : check this part supposed
	// to be userCollection.UpdateMan... and also
	// remove productCol.. from the function
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil

}
func BuyItemFromCart(
	ctx context.Context,
	userCollection *mongo.Collection,
	userId string) error {
	//	fetch the cart of the user
	//	find the cart total
	//	create an order with the items
	// added order to the user collection
	// added items in the cart of the order list
	//	empty up the cart
	_user_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return ErrUserIdIsNotValid
	}
	var getCartItems models.User
	var orderCart models.Order

	orderCart.OrderId = primitive.NewObjectID()
	orderCart.OrderedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderCart.OrderCart = make([]models.ProductUser, 0)
	orderCart.PaymentMethod.COD = true

	unwind := bson.D{
		{Key: "$unwind", Value: bson.D{
			primitive.E{Key: "path", Value: "$user_cart"},
		}}}
	grouping := bson.D{
		{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: "$_id"},
			{Key: "total", Value: bson.D{
				primitive.E{Key: "$sum", Value: "user_cart.price"},
			}},
		}},
	}
	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		unwind, grouping,
	})
	ctx.Done()
	if err != nil {
		panic(err)
	}
	var getUserCart []bson.M
	err = currentResults.All(ctx, &getUserCart)
	if err != nil {
		panic(err)
		return err
	}
	var total_price int32
	for _, user_item := range getUserCart {
		price := user_item["total"]
		total_price = price.(int32)
	}
	orderCart.Price = int(total_price)
	filter := bson.D{
		primitive.E{Key: "_id", Value: _user_id},
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			primitive.E{Key: "orders", Value: orderCart},
		}},
	}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return err
	}
	err = userCollection.FindOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: _user_id},
	}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
		return err
	}
	filter2 := bson.D{
		primitive.E{Key: "_id", Value: _user_id},
	}
	update2 := bson.M{
		"$push": bson.M{"orders.$[].order_list": bson.M{
			"$each": getCartItems.UserCart,
		}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	user_cart_empty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: _user_id}}
	update3 := bson.D{{
		Key: "$set", Value: bson.D{
			primitive.E{Key: "user_cart", Value: user_cart_empty},
		}}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}
	return nil

}

func InstantBuy(
	ctx context.Context,
	prodCollection *mongo.Collection,
	userCollection *mongo.Collection,
	productId primitive.ObjectID,
	userId string) error {
	_user_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var product_details models.ProductUser
	var orders_detail models.Order

	orders_detail.OrderId = primitive.NewObjectID()
	orders_detail.OrderedAt = time.Now()
	orders_detail.OrderCart = make([]models.ProductUser, 0)
	orders_detail.PaymentMethod.COD = true

	err = prodCollection.FindOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: productId},
	}).Decode(&product_details)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}
	orders_detail.Price = product_details.Price

	filter := bson.D{
		primitive.E{Key: "_id", Value: _user_id},
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			primitive.E{Key: "orders", Value: orders_detail},
		}},
	}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	filter2 := bson.D{
		primitive.E{Key: "_id", Value: _user_id},
	}
	update2 := bson.M{"$push": bson.M{
		"orders.$[].order_list": product_details,
	}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil
}

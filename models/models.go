package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName      *string            `json:"first_name" validate:"required,min=2,max=100" bson:"first_name"`
	LastName       *string            `json:"last_name" validate:"required,min=2,max=100" bson:"last_name"`
	Email          *string            `json:"email" validate:"required" bson:"email"`
	Password       *string            `json:"password" validate:"required,min=5,max=100" bson:"password"`
	Phone          *string            `json:"phone" validate:"required,min=10,max=20" bson:"phone"`
	Token          *string            `json:"token" bson:"token"`
	RefreshToken   *string            `json:"refresh_token" bson:"refresh_token"`
	CreatedAT      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAT      time.Time          `json:"updated_at" bson:"updated_at"`
	UserId         string             `json:"user_id" bson:"user_id"`
	UserCart       []ProductUser      `json:"user_cart" bson:"user_cart"`
	AddressDetails []Address          `json:"address" bson:"address"`
	OrderStatus    []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	ProductId   primitive.ObjectID `bson:"_id"`
	ProductName *string            `json:"product_name" bson:"product_name"`
	Price       *uint64            `json:"price" bson:"price"`
	Rating      *uint8             `json:"rating" bson:"rating"`
	Image       *string            `json:"image" bson:"image"`
}
type ProductUser struct {
	ProductId   primitive.ObjectID `bson:"_id"`
	ProductName *string            `json:"product_name" bson:"product_name"`
	Price       int                `json:"price" bson:"price"`
	Rating      *uint8             `json:"rating" bson:"rating"`
	Image       *string            `json:"image" bson:"image"`
}

type Address struct {
	AddressId primitive.ObjectID `bson:"_id"`
	House     *string            `json:"house_name" bson:"house_name"`
	Street    *string            `json:"street_name" bson:"street_name"`
	City      *string            `json:"city_name" bson:"city_name"`
	PinCode   *string            `json:"pin_code" bson:"pin_code"`
}

type Order struct {
	OrderId       primitive.ObjectID `bson:"_id"`
	OrderCart     []ProductUser      `json:"order_cart" bson:"order_cart"`
	OrderedAt     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price         int                `json:"price" bson:"price"`
	Discount      *int               `json:"discount" bson:"discount"`
	PaymentMethod Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" validate:"eq|true|false" bson:"digital"`
	COD     bool `json:"cod" validate:"eq|true|false" bson:"cod"`
}

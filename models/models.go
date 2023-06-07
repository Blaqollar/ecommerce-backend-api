package models

import "time"

type User struct {
	ID             string        `json:"_id" bson:"_id,omitempty"`
	FirstName      string        `json:"first_name" bson:"first_name,omitempty" validate:"required, min=2, max=30"`
	LastName       string        `json:"last_name" bson:"last_name,omitempty" validate:"required, min=2, max=30"`
	Password       string        `json:"password" bson:"password,omitempty" validate:"required, min=6"`
	Email          string        `json:"email" bson:"email,omitempty" validate:"required"`
	Phone          string        `json:"phone" bson:"phone,omitempty" validate:"required"`
	Token          string        `json:"token" bson:"token,omitempty"`
	RefreshToken   string        `json:"refresh_token" bson:"refresh_token,omitempty"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at,omitempty"`
	UserId         string        `json:"user_id" bson:"user_id,omitempty"`
	UserCart       []ProductUser `json:"user_cart" bson:"user_cart,omitempty"`
	AddressDetails []Address     `json:"address_details" bson:"address_details,omitempty"`
	OrderStatus    []Order       `json:"order_status" bson:"order_status,omitempty"`
}

type Product struct {
	ProductID   string  `json:"product_id" bson:"product_id,omitempty"`
	ProductName string  `json:"product_name" bson:"product_name,omitempty"`
	Description string  `json:"description" bson:"description,omitempty"`
	Price       float64 `json:"price" bson:"price,omitempty"`
	Rating      float64 `json:"rating" bson:"rating,omitempty"`
	Image       string  `json:"image" bson:"image,omitempty"`
}

type ProductUser struct {
	ProductID   string  `json:"product_id" bson:"product_id,omitempty"`
	ProductName string  `json:"product_name" bson:"product_name"`
	Price       float64 `json:"price" bson:"price,omitempty"`
	Rating      float64 `json:"rating" bson:"rating,omitempty"`
	Image       string  `json:"image" bson:"image,omitempty"`
}

type Address struct {
	AddressID string `json:"address_id" bson:"address_id,omitempty"`
	House     string `json:"house_name" bson:"house_name,omitempty"`
	Street    string `json:"street_name" bson:"street_name,omitempty"`
	City      string `json:"city_name" bson:"city_name,omitempty"`
	Pincode   string `json:"pincode" bson:"pincode,omitempty"`
}

type Order struct {
	OrderId       string        `json:"orderId" bson:"order,omitempty"`
	OrderCart     []ProductUser `json:"order_list" bson:"order_list,omitempty"`
	OrderedAt     time.Time     `json:"ordered_at" bson:"ordered_at,omitempty"`
	Price         float64       `json:"price" bson:"price,omitempty"`
	Discount      float64       `json:"discount" bson:"discount,omitempty"`
	PaymentMethod Payment       `json:"payment_method" bson:"payment_method,omitempty"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital,omitempty"`
	COD     bool `json:"cash_on_delivery" bson:"cash_on_delivery,omitempty"`
}

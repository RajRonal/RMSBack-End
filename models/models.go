package models

import (
	"time"
)

type CreateUser struct {
	Name     string `json:"name" db:"user_name"`
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type UserRole string

type ClaimsKey string

const (
	UserRoleAdmin    UserRole  = "admin"
	UserRoleSubAdmin UserRole  = "sub-admin"
	UserRoleUser     UserRole  = "user"
	ClaimKey         ClaimsKey = "claim"
)

type CreateRole struct {
	ID       string   `db:"id" json:"userId"`
	Role     UserRole `json:"role" db:"user_role"`
	Username string   `json:"username" db:"username"`
}
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Session struct {
	SessionID string    `db:"session_id" json:"sessionId"`
	ExpiredAt time.Time `db:"expired_at" json:"expiredAt"`
	ID        string    `db:"id" json:"id"`
}
type AddLogin struct {
	ID       string `db:"id" json:"id"`
	Password string `db:"password" json:"password"`
}
type GetRole struct {
	Role UserRole `json:"role" db:"user_role"`
}
type Restaurant struct {
	ID         string  `db:"restaurant_id" json:"id"`
	Name       string  `db:"restaurant_name" json:"name"`
	Longitude  float64 `db:"longitude" json:"longitude"`
	Latitude   float64 `db:"latitude" json:"latitude"`
	TotalCount int     `json:"-" db:"total_count"`
}
type PaginatedTask struct {
	Data       []PaginatedData `json:"data"`
	TotalCount int             `json:"totalCount"`
}

type PaginatedData struct {
	Name       string `json:"name" db:"user_name"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password"`
	TotalCount int    `json:"-" db:"total_count"`
}

type PaginatedRestaurant struct {
	Data       []Restaurant `json:"data"`
	TotalCount int          `json:"totalCount"`
}

type Dishes struct {
	DishName   string  `json:"dishName" db:"dish_name"`
	DishPrice  float64 `json:"dishPrice" db:"dish_price"`
	TotalCount int     `json:"-" db:"total_count"`
}

type PaginatedDishes struct {
	Data       []Dish `json:"data"`
	TotalCount int    `json:"totalCount"`
}

type FetchDish struct {
	RestaurantID string `json:"restaurantId" db:"restaurant_id"`
}

type Dish struct {
	DishID     string  `json:"dishId"  db:"dish_id"`
	DishName   string  `json:"dishName" db:"dish_name"`
	DishPrice  float64 `json:"dishPrice" db:"dish_price"`
	TotalCount int     `json:"-" db:"total_count"`
}

type Location struct {
	Longitude float64 `json:"longitude"  db:"longitude"`
	Latitude  float64 `json:"latitude"  db:"latitude"`
}

type LocationDistance struct {
	Distance float64 `json:"distance"`
}

type UpdateRestaurant struct {
	Name string `db:"restaurant_name" json:"name"`
}

type UpdateDish struct {
	DishName string `db:"dish_name" json:"dishName"`
}

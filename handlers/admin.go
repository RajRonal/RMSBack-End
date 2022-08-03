package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"

	"Rms/database/helper"
	"Rms/models"
)

func FetchAllUserRole(writer http.ResponseWriter, request *http.Request) {
	pageNo := request.URL.Query().Get("page")
	if pageNo == "" {
		pageNo = "0"
	}
	Page, err := strconv.Atoi(pageNo)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	limitSize := request.URL.Query().Get("limit")
	if limitSize == "" {
		limitSize = "5"
	}

	Limit, err := strconv.Atoi(limitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userRole := request.URL.Query().Get("role")
	status := make(pq.StringArray, 0)
	var isRole bool
	if string(models.UserRoleUser) == userRole || string(models.UserRoleSubAdmin) == userRole {
		isRole = true
		status = append(status, userRole)
	} else {
		status = append(status, string(models.UserRoleUser), string(models.UserRoleSubAdmin))
		// status = append(status, "sub-admin")
	}

	user, err := helper.GetAllUserData(status, Page, Limit, isRole)
	if err != nil {
		logrus.Error("FetchAllRole: Sub Admin can't be fetched %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := json.NewEncoder(writer).Encode(user)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func ChangeUserRole(writer http.ResponseWriter, request *http.Request) {
	UserID := chi.URLParam(request, "userId")
	err := helper.ChangeRole(UserID)
	if err != nil {
		logrus.Error("ChangeUserRole: User role can't be fetched %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	printMessage := "User Role changed Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateRestaurant(writer http.ResponseWriter, request *http.Request) {
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var restaurant models.Restaurant
	err := json.NewDecoder(request.Body).Decode(&restaurant)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = helper.AddRestaurant(restaurant.Name, uc.ID, restaurant.Longitude, restaurant.Latitude)
	if err != nil {
		logrus.Error(" FeedIntoRestaurant:error in adding restaurant %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := " Restaurant Created Successfully"
	errs := json.NewEncoder(writer).Encode(message)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func FetchAllRestaurant(writer http.ResponseWriter, request *http.Request) {
	searchRestaurant := request.URL.Query().Get("search")
	pageNo := request.URL.Query().Get("page")
	if pageNo == "" {
		pageNo = "0"
	}
	Page, err := strconv.Atoi(pageNo)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	limitSize := request.URL.Query().Get("limit")
	if limitSize == "" {
		limitSize = "5"
	}
	Limit, err := strconv.Atoi(limitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	restaurant, err := helper.GetAllRestaurant(Page, Limit, searchRestaurant)
	if err != nil {
		logrus.Error("FetchAllRestaurant: restaurant  can't be fetched %V", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := json.NewEncoder(writer).Encode(restaurant)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func FetchAllDish(writer http.ResponseWriter, request *http.Request) {
	RestaurantID := chi.URLParam(request, "id")
	pageNo := request.URL.Query().Get("page")
	if pageNo == "" {
		pageNo = "0"
	}
	Page, err := strconv.Atoi(pageNo)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	limitSize := request.URL.Query().Get("limit")
	if limitSize == "" {
		limitSize = "5"
	}
	Limit, err := strconv.Atoi(limitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	dishes, err := helper.GetAllDishes(RestaurantID, Page, Limit)
	if err != nil {
		logrus.Error("FetchAllDish: restaurant  can't be fetched %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := json.NewEncoder(writer).Encode(dishes)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateRestaurant(writer http.ResponseWriter, request *http.Request) {
	var restaurant models.UpdateRestaurant
	RestaurantID := chi.URLParam(request, "id")
	err := json.NewDecoder(request.Body).Decode(&restaurant)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = helper.UpdateRestaurant(RestaurantID, restaurant.Name)
	if err != nil {
		logrus.Error("UpdateRestaurant: can't update restaurant %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	printMessage := "Restaurant updated Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteRestaurant(writer http.ResponseWriter, request *http.Request) {
	RestaurantID := chi.URLParam(request, "id")
	err := helper.DeleteRestaurant(RestaurantID)
	if err != nil {
		logrus.Error("DeleteRestaurant: can't delete  %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	printMessage := "Restaurant Deleted Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateDish(writer http.ResponseWriter, request *http.Request) {
	RestaurantID := chi.URLParam(request, "id")
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var dish models.Dishes
	err := json.NewDecoder(request.Body).Decode(&dish)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = helper.InsertDish(dish.DishName, RestaurantID, uc.ID, dish.DishPrice)
	if err != nil {
		logrus.Error("FeedIntoDishes:error in adding dishes %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	printMessage := "Dishes Created Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteDish(writer http.ResponseWriter, request *http.Request) {
	DishID := chi.URLParam(request, "dishId")
	err := helper.DeleteDish(DishID)
	if err != nil {
		logrus.Error("DeleteDish: can't delete dish %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	printMessage := "Dish Deleted Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateDish(writer http.ResponseWriter, request *http.Request) {
	DishID := chi.URLParam(request, "dishId")
	var user models.UpdateDish
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = helper.UpdateDish(DishID, user.DishName)
	if err != nil {
		logrus.Error("update Error : can't update dish %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	printMessage := "Dish updated Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

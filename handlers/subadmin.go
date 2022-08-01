package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"Rms/database/helper"
	"Rms/models"
)

func FetchAllRestaurantSubAdmin(writer http.ResponseWriter, request *http.Request) {
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
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := helper.GetAllRestaurantSubAdmin(uc.ID, Page, Limit)
	if err != nil {
		logrus.Error("FetchAllRestaurantSubAdmin: restaurant  can't be fetched")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := json.NewEncoder(writer).Encode(user)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteRestaurantSubAdmin(writer http.ResponseWriter, request *http.Request) {
	RestaurantID := chi.URLParam(request, "id")
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resID := helper.CheckUpdate(RestaurantID)
	if uc.ID != resID {
		logrus.Error("DeleteRestaurantSubAdmin: can't delete dish")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := helper.DeleteRestaurantSubAdmin(RestaurantID, uc.ID)
	if err != nil {
		logrus.Error("DeleteRestaurantSubAdmin : can't delete restaurant")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	printMessage := "Restaurant Deleted Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateRestaurantSubAdmin(writer http.ResponseWriter, request *http.Request) {
	var restaurant models.UpdateRestaurant
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(request.Body).Decode(&restaurant)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	RestaurantID := chi.URLParam(request, "id")
	checkUpdate := helper.CheckUpdate(RestaurantID)
	if checkUpdate != uc.ID {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = helper.UpdateRestaurantSubAdmin(RestaurantID, restaurant.Name, uc.ID)
	if err != nil {
		logrus.Error("UpdateRestaurantSubAdmin: can't update error")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	printMessage := "Restaurant updated Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func FetchAllDishSubAdmin(writer http.ResponseWriter, request *http.Request) {
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
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	RestaurantID := chi.URLParam(request, "id")
	dish, err := helper.GetAllDishesSubAdmin(RestaurantID, uc.ID, Page, Limit)
	if err != nil {
		logrus.Error("FetchAllDishSubAdmin: restaurant  can't be fetched")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := json.NewEncoder(writer).Encode(dish)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteDishSubAdmin(writer http.ResponseWriter, request *http.Request) {
	DishID := chi.URLParam(request, "dishId")
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	check := helper.CheckDeleteDish(DishID)
	if check != uc.ID {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := helper.DeleteDishSubAdmin(DishID, uc.ID)
	if err != nil {
		logrus.Error(" DeleteDishSubAdmin: can't delete dish")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	printMessage := "Dish Deleted Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateDishSubAdmin(writer http.ResponseWriter, request *http.Request) {
	DishID := chi.URLParam(request, "dishId")
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.UpdateDish
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	check := helper.CheckDeleteDish(DishID)
	if check != uc.ID {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = helper.UpdateDishSubAdmin(DishID, user.DishName)
	if err != nil {
		logrus.Error("UpdateDishSubAdmin: can't Update dish")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	printMessage := "Dish updated Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

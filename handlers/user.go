package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"

	"Rms/claims"
	"Rms/database/helper"
	"Rms/models"
)

//todo: use logrus
var JwtKey = []byte("secureSecretText")

func SignUp(writer http.ResponseWriter, request *http.Request) {
	var user models.CreateUser

	error := json.NewDecoder(request.Body).Decode(&user)
	if error != nil {
		logrus.Error("SignUp: Error in decoding json %v", error)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := helper.CreateUser(user)
	if err != nil {
		logrus.Error("SignUp: Error in Creating user %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = helper.CreateRole(userID, user.Username, string(models.UserRoleUser))
	if err != nil {
		logrus.Error("SignUp:error in allocating role to the user %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	message := "User Creation successful"
	errs := json.NewEncoder(writer).Encode(message)
	if errs != nil {
		logrus.Error("SignUp: Error in Encoding json %v", errs)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credential
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		logrus.Error("Login: Error in decoding json %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	passkey, err := helper.LoginUser(credentials.Username)
	if err != nil {
		logrus.Error("Login: Error in getting password %v", err)
		w.WriteHeader(http.StatusUnauthorized)
	}

	if passkey.Password != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
	}

	userRole, err := helper.GetRole(credentials.Username)
	if err != nil {
		logrus.Error("Login: Error in Getting Role of User %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	expirationTime := time.Now().Add(time.Hour * 7)
	sessionID, err := helper.CreateSession(passkey.ID, expirationTime)
	if err != nil {
		logrus.Error("LogIn : Error in Creating the session %v", err)
		w.WriteHeader(http.StatusUnauthorized)
	}

	mapClaim := &claims.MapClaims{
		SessionID: sessionID,
		ID:        passkey.ID,
		Role:      userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaim)
	signedToken, err := token.SignedString([]byte("secureSecretText"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tokenByte, err := json.Marshal(signedToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(signedToken)
	_, _ = w.Write(tokenByte)
}

func SetLocation(writer http.ResponseWriter, request *http.Request) {
	var userLocation models.Location
	err := json.NewDecoder(request.Body).Decode(&userLocation)
	if err != nil {
		logrus.Error("SetLocation: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = helper.SetLocation(userLocation.Longitude, userLocation.Latitude, uc.ID)
	if err != nil {
		logrus.Error("SetLocation: Error in Feeding The location %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	printMessage := "Location Created Successfully"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
func GetDistance(writer http.ResponseWriter, request *http.Request) {
	restaurantID := chi.URLParam(request, "restaurantId")
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	RestaurantDistance, err := helper.GetRestaurantLongitudeLatitude(restaurantID)
	if err != nil {
		logrus.Error("GetDistance: Error in Getting Restaurant Longitude and Latitude %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	UserDistance, err := helper.GetUserLongitudeLatitude(uc.ID)
	if err != nil {
		logrus.Error("GetDistance: Error in Getting User Longitude and Latitude %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	Distance, err := helper.GetDistance(RestaurantDistance.Longitude, RestaurantDistance.Latitude, UserDistance.Longitude, UserDistance.Latitude)
	if err != nil {
		logrus.Error("GetDistance: Error in Getting Distance %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := json.NewEncoder(writer).Encode(Distance)
	if errs != nil {
		log.Printf("GetDistance: Error in Encoding details %v", errs)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func SignOut(writer http.ResponseWriter, request *http.Request) {
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	deleteSession := helper.DeleteSession(uc.SessionID)
	if deleteSession != nil {
		logrus.Error("SignOut: Session can't be deleted %v", deleteSession)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	printMessage := "Sign Out Successful"
	errs := json.NewEncoder(writer).Encode(printMessage)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

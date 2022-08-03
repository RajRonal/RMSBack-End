package helper

import (
	"Rms/claims"
	"Rms/database"
	"Rms/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
)

func CreateUser(user models.CreateUser) (string, error) {
	SQL := `INSERT INTO users (user_name, email, username, password)
			VALUES($1, $2, $3, $4)
			returning id`
	var userID string
	err := database.DB.Get(&userID, SQL, user.Name, user.Email, user.Username, user.Password)
	if err != nil {
		logrus.Error("CreateUser:Error in creating user %v", err)
		return "", err
	}
	return userID, nil
}

func CreateRole(id, username, role string, tx *sqlx.Tx) error {
	SQL := `INSERT INTO roles(user_id,user_role,username)
           VALUES($1,$2,$3)`
	_, err := database.DB.Exec(SQL, id, role, username)
	if err != nil {
		logrus.Error("CreateRole: error occurred in creating role %v", err)
		return err
	}
	return nil
}
func CreateSession(id string, expiredAt time.Time) (string, error) {
	SQL := `INSERT INTO sessions (id,expired_at ) 	
			VALUES ($1, $2)  
			returning session_id;`
	var sessID string
	err := database.DB.Get(&sessID, SQL, id, expiredAt)
	if err != nil {
		logrus.Error("CreateSession:Error in creating Session %v", err)
		return "", err
	}
	return sessID, nil
}

func SessionExist(sessionID string) (bool, error) {
	var isExpired bool
	query := `SELECT count(*) > 0
			  FROM sessions 
			WHERE session_id=$1  and expired_at >now() and archived_at is null`
	checkSessionErr := database.DB.Get(&isExpired, query, sessionID)
	if checkSessionErr != nil {
		logrus.Error("SessionExist:Error in checking existence of session %v", checkSessionErr)
		return isExpired, checkSessionErr
	}
	return isExpired, nil
}
func LoginUser(username string) (*models.AddLogin, error) {
	SQL := `SElECT id,password from users where username=$1`
	var pass models.AddLogin
	err := database.DB.Get(&pass, SQL, username)
	if err != nil && err != sql.ErrNoRows {
		logrus.Error("LoginUser:Error in logging in %v", err)
		return nil, err
	}
	return &pass, nil
}
func GetRole(username string) (string, error) {
	SQL := `SElECT user_role from roles where username=$1`
	var pass string
	err := database.DB.Get(&pass, SQL, username)
	if err != nil {
		logrus.Error("GetRole: error in fetching data: %v", err)
		return "", err
	}
	return pass, nil
}

func AddRestaurant(name, id string, long, lati float64) (string, error) {
	SQL := `INSERT INTO restaurant (restaurant_name, created_by, longitude, latitude ) 	
			VALUES ($1, $2, $3, $4)`
	var user string
	_, err := database.DB.Exec(SQL, name, id, long, lati)
	if err != nil {
		log.Printf("AddRestaurant: error in adding restaurant %v", err)
		return "", err
	}
	return user, nil
}

func GetAllUserData(userRole pq.StringArray, pageNo, taskSize int, isRole bool) (models.PaginatedTask, error) {
	var data models.PaginatedTask
	SQL := `WITH userTask AS (SELECT count(*) over () total_count, user_name, email, password
							  FROM users
							  JOIN roles ON users.id = roles.user_id
							  WHERE (($4 OR user_role= ANY($3)))
								AND roles.archived_at IS NULL)
			SELECT total_count, user_name, email, password
			from userTask
			LIMIT $1 OFFSET $2
			`
	user := make([]models.PaginatedData, 0)
	err := database.DB.Select(&user, SQL, taskSize, pageNo*taskSize, userRole, !isRole)
	if err != nil {
		logrus.Error("GetAllData: error in fetching data: %v", err)
		return data, err
	}
	if len(user) == 0 {
		return data, err
	}
	data.TotalCount = user[0].TotalCount
	data.Data = user
	return data, err
}

func GetAllRestaurant(pageNo, taskSize int, search string) (models.PaginatedRestaurant, error) {
	var data models.PaginatedRestaurant
	SQL := `WITH getRestaurant AS (SELECT  count(*) over ()total_count, restaurant_id,restaurant_name, longitude, latitude
			FROM restaurant where restaurant.restaurant_name ILIKE '%' || $3 ||'%' AND  archived_at IS NULL)
			
			 SELECT  total_count,restaurant_id,restaurant_name, longitude, latitude from getRestaurant 
			        LIMIT $1
					OFFSET $2`
	user := make([]models.Restaurant, 0)
	err := database.DB.Select(&user, SQL, taskSize, pageNo*taskSize, search)
	if err != nil {
		logrus.Error("GetAllRestaurant: error in fetching data: %v", err)
		return data, err
	}
	if len(user) == 0 {
		return data, err
	}
	data.TotalCount = user[0].TotalCount
	data.Data = user
	return data, err
}
func InsertDish(dishname, restauarntID, userID string, dishPrice float64) (string, error) {
	SQL := `INSERT INTO dishes (dish_name, dish_price, restaurant_id, created_by)
	VALUES ($1, $2, $3, $4)`
	var user string
	_, err := database.DB.Exec(SQL, dishname, dishPrice, restauarntID, userID)
	if err != nil {
		log.Printf("InsertDishes: Error in Asdding Dishes %v", err)
		return "", err
	}
	return user, nil
}
func GetAllDishes(id string, pageNo, taskSize int) (models.PaginatedDishes, error) {
	var data models.PaginatedDishes
	SQL := `WITH getDishes AS (SELECT  count(*) over ()total_count,dish_id,dish_name,dish_price
			FROM  dishes where restaurant_id=$3 and archived_at IS NULL)
			
			 SELECT  total_count,dish_id,dish_name,dish_price from getDishes 
			        LIMIT $1
					OFFSET $2`
	user := make([]models.Dish, 0)
	err := database.DB.Select(&user, SQL, taskSize, pageNo*taskSize, id)
	if err != nil {
		log.Printf("GetAllDishes: error in fetching data: %v", err)
		return data, err
	}
	if len(user) == 0 {
		return data, err
	}
	data.TotalCount = user[0].TotalCount
	data.Data = user
	return data, err
}
func ChangeRole(userId string) error {
	SQL := `UPDATE roles
	SET user_role=$1
	WHERE user_id=$2 and archived_at is null`
	_, err := database.DB.Exec(SQL, models.UserRoleSubAdmin, userId)
	if err != nil {
		logrus.Error("ChangeRole: Error in changing role %v", err)
		return err
	}
	return nil
}
func DeleteDish(id string) error {
	SQL := `UPDATE dishes
	SET archived_at = now()
	WHERE dish_id = $1 and archived_at is null`
	_, err := database.DB.Exec(SQL, id)
	if err != nil {
		logrus.Error("DeleteDish:error in deleting dish %v", err)
		return err
	}
	return nil
}
func UpdateDish(dishID, name string) error {
	SQL := `UPDATE dishes
	SET dish_name=$2
	WHERE dish_id=$1 and archived_at is null`
	_, err := database.DB.Exec(SQL, dishID, name)
	if err != nil {
		logrus.Error("UpdateDish: Error in updating dish %v", err)
		return err
	}
	return nil
}
func DeleteRestaurant(restaurantId string) error {
	SQL := `UPDATE restaurant
	SET archived_at=now()
	WHERE restaurant_id=$1 and archived_at is null`
	_, err := database.DB.Exec(SQL, restaurantId)
	if err != nil {
		logrus.Error("DeleteRestaurant:Error in deleting dish %v", err)
		return err
	}
	return nil
}
func GetAllRestaurantSubAdmin(id string, pageNo, taskSize int) (models.PaginatedRestaurant, error) {
	var data models.PaginatedRestaurant
	SQL := `WITH getRestaurant AS (SELECT  count(*) over ()total_count, restaurant_id,restaurant_name, longitude, latitude
			FROM restaurant WHERE created_by=$3 and archived_at IS NULL)
			
			 SELECT  total_count,restaurant_id,restaurant_name, longitude, latitude from getRestaurant 
			        LIMIT $1
					OFFSET $2`
	user := make([]models.Restaurant, 0)
	err := database.DB.Select(&user, SQL, taskSize, pageNo*taskSize, id)
	if err != nil {
		logrus.Error("GetAllRestaurantSubAdmin: error in fetching data: %v", err)
	}
	if len(user) == 0 {
		return data, err
	}
	data.TotalCount = user[0].TotalCount
	data.Data = user
	return data, err
}
func GetAllDishesSubAdmin(id, userid string, pageNo, taskSize int) (models.PaginatedDishes, error) {
	var data models.PaginatedDishes
	SQL := `WITH getDishes AS (SELECT  count(*) over ()total_count,dish_id,dish_name,dish_price
			FROM  dishes where restaurant_id=$3 and created_by=$4 and archived_at IS NULL )
			
			 SELECT  total_count,dish_id,dish_name,dish_price from getDishes 
			        LIMIT $1
					OFFSET $2`
	user := make([]models.Dish, 0)
	err := database.DB.Select(&user, SQL, taskSize, pageNo*taskSize, id, userid)
	if err != nil {
		logrus.Error("GetAllDishesSubAdmin: error in fetching data: %v", err)
		return data, err
	}
	if len(user) == 0 {
		return data, err
	}
	data.TotalCount = user[0].TotalCount
	data.Data = user
	return data, err
}
func DeleteDishSubAdmin(dishId, userid string) error {
	SQL := `UPDATE dishes
	SET archived_at=now()
	WHERE dish_id=$1 and created_by=$2 and archived_at is null`
	_, err := database.DB.Exec(SQL, dishId, userid)
	if err != nil {
		logrus.Error("DeleteDishSubAdmin: Error in deleting dish subadmin %v", err)
		return err
	}
	return nil
}
func UpdateDishSubAdmin(dishId, name string) error {
	SQL := `UPDATE dishes
	SET dish_name=$2
	WHERE dish_id=$1  and archived_at is Null`
	_, err := database.DB.Exec(SQL, dishId, name)
	if err != nil {
		logrus.Error("UpdateDishSubAdmin: Error in uodating dish in subadmin role %v", err)
		return err
	}
	return nil
}
func DeleteRestaurantSubAdmin(restaurantId, userID string) error {
	SQL := `UPDATE restaurant
	SET archived_at=now()
	WHERE restaurant_id=$1 and created_by=$2 and archived_at is NULL`
	_, err := database.DB.Exec(SQL, restaurantId, userID)
	if err != nil {
		logrus.Error("DeleteRestaurantSubAdmin: Error in deleting restaurant in subadmin role %v", err)
		return err
	}
	return nil
}
func SetLocation(longitude, latitude float64, userid string) (string, error) {
	SQL := `INSERT INTO location (longitude,latitude,user_id ) 	
			VALUES ($1, $2,$3)  
			`
	var user string
	_, err := database.DB.Exec(SQL, longitude, latitude, userid)
	if err != nil {
		logrus.Error("SetLocation: Error in feeding location %v", err)
		return "", err
	}
	return user, nil
}
func GetRestaurantLongitudeLatitude(restaurantId string) (models.Location, error) {
	var user models.Location
	SQL := `SELECT longitude,latitude FROM restaurant WHERE restaurant_id=$1 and archived_at is null`
	err := database.DB.Get(&user, SQL, restaurantId)
	if err != nil {
		logrus.Error("GetRestaurantLongitudeLatitude: error in fetching restautant distance %v", err)
		return user, err
	}
	return user, nil
}

func GetUserLongitudeLatitude(id string) (models.Location, error) {
	var user models.Location
	SQL := `SELECT longitude,latitude FROM location WHERE user_id=$1`
	err := database.DB.Get(&user, SQL, id)
	if err != nil {
		logrus.Error("GetUserLongitudeLatitude: Error in getting user distance %v", err)
		return user, err
	}
	return user, nil

}

func GetDistance(restaurantLongitude, restaurantLatitude, userLongitude, userLatitude float64) (models.LocationDistance, error) {
	var result models.LocationDistance
	SQL := `SELECT (point($1,$2) <-> point($3,$4))as distance`
	err := database.DB.Get(&result, SQL, userLongitude, userLatitude, restaurantLongitude, restaurantLatitude)
	if err != nil {
		logrus.Error("GetLocation: Error in calculating distance %v", err)
		return result, err
	}
	return result, nil
}
func UpdateRestaurant(restaurantId, restaurantName string) error {
	SQL := `UPDATE restaurant
	SET restaurant_name=$1
	WHERE restaurant_id=$2 AND archived_at IS NULL`
	_, err := database.DB.Exec(SQL, restaurantName, restaurantId)
	if err != nil {
		logrus.Error("UpdateRestaurant: error in updating restaurant %v", err)
		return err
	}
	return nil
}
func UpdateRestaurantSubAdmin(restaurantId, restaurantName, userID string) error {
	SQL := `UPDATE restaurant
	SET restaurant_name=$1
	WHERE restaurant_id=$2 AND archived_at IS NULL and created_by=$3`
	_, err := database.DB.Exec(SQL, restaurantName, restaurantId, userID)
	if err != nil {
		logrus.Error("UpdateRestaurantSubAdmin: Error in updating restaurant in sub-admin role %v", err)
		return err
	}
	return nil
}
func CheckUpdate(resID string) string {
	var user string
	SQL := `SELECT created_by from restaurant 
			where restaurant_id=$1 `
	err := database.DB.Get(&user, SQL, resID)
	if err != nil {
		logrus.Error("CheckUpdate: Error in checking update%v", err)
		return err.Error()
	}
	return user
}
func CheckDeleteDish(restaurantId string) string {
	var user string
	SQL := `SELECT created_by from dishes 
			where dish_id=$1 `
	err := database.DB.Get(&user, SQL, restaurantId)
	if err != nil {
		logrus.Error("CheckDeleteDish:Error in checking dish owner %v", err)
		return err.Error()
	}
	return user
}
func DeleteSession(sessionID string) error {
	currentTime := time.Now()
	SQL := `UPDATE sessions
			  SET archived_at= $1,
			      expired_at= now()
			  WHERE session_id= $2`
	_, err := database.DB.Exec(SQL, currentTime, sessionID)
	if err != nil {
		logrus.Error("DeleteSession: error in Deleting Session: %v", err)
		return err
	}
	return err
}
func GetContextData(request *http.Request) *claims.MapClaims {
	uc, ok := request.Context().Value(models.ClaimKey).(*claims.MapClaims)
	if !ok {
		logrus.Error("ContextData: Error In Parsing Context")
		return nil
	}
	return uc
}

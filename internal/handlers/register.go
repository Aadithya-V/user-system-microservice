package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Aadithya-V/user-system-microservice/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *redis.Client) func(ctx *gin.Context) {
	fx := func(ctx *gin.Context) {
		var newuser models.User
		if err := ctx.BindJSON(&newuser); err != nil {
			ctx.JSON(http.StatusBadRequest, &gin.H{"error": "registration information not submitted correctly."}) // print json format for reference
			return
		}
		// TODO:validate client data? Better done at the client..

		if err := db.HGet(CTX, "users", newuser.Name).Err(); err != redis.Nil {
			ctx.JSON(http.StatusConflict, &gin.H{"error": "Username not available. Already in use."}) //Do you allow reuse of deleted accounts' usernames?
			return
		}

		id, err := db.Incr(CTX, "next_user_id").Result()
		if err == redis.Nil { //err if next_user_id not init in db
			ctx.JSON(http.StatusInternalServerError, &gin.H{"error": "unable to create id. DB not initialised."})
			return
		}
		newuser.ID = strconv.FormatInt(id, 10)

		db.HSet(CTX, "users", newuser.Name, newuser.ID)

		hash, _ := bcrypt.GenerateFromPassword([]byte(newuser.Pwd), bcrypt.DefaultCost)
		newuser.Pwd = string(hash)

		db.GeoAdd(CTX, "locations", &redis.GeoLocation{
			Longitude: newuser.Longitude,
			Latitude:  newuser.Latitude,
			Name:      newuser.ID,
		})

		db.HSet(CTX, "user:"+newuser.ID, "id", newuser.ID, "name", newuser.Name, "dob", newuser.DOB, "address", newuser.Address, "latitude", newuser.Latitude, "longitude", newuser.Longitude, "description", newuser.Description, "createdAt", time.Now(), "pwd", newuser.Pwd)

		ctx.JSON(http.StatusCreated, &gin.H{"message": "account successfully created."})
	}
	return fx
}

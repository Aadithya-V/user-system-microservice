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
			//ctx.JSON(http., struct{}{})
			return
		}
		// validate client data

		if res := db.Get(CTX, newuser.Name); res.Err() != redis.Nil {
			return //error status and info
		}

		id := db.Incr(CTX, "next_user_id").Val() //err if next_user_id not init in db
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

		ctx.JSON(http.StatusCreated, struct{}{})

	}
	return fx
}

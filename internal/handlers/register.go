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
		lat, lon := newuser.Coordinates[0], newuser.Coordinates[1]
		if err := validateCoordinates(lat, lon); err != nil || newuser.Coordinates == NULLISLAND {
			ctx.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()}) // Received Invalid Coordinates
			return
		}

		if err := db.HGet(CTX, "users", newuser.Name).Err(); err != redis.Nil {
			ctx.JSON(http.StatusConflict, &gin.H{"error": "Username not available. Already in use."}) //Do you allow reuse of deleted accounts' usernames?
			return
		}

		id, err := db.Incr(CTX, "next_user_id").Result()
		if err == redis.Nil {
			ctx.JSON(http.StatusInternalServerError, &gin.H{"error": "unable to create id."})
			return
		} /* 	// err != nil cannot ever happen practically (only if size of id string value exceeds 512 MiB).
		// If next_user_id not initialized in db to desired value, incr starts by returning 1.
		// Collission will happen if key is mistakenly deleted as 1 will be returned again. All further registrations
		// thence will be erraneous. So, (is it?) worth testing for only on case ID==1 if an id=1 already exists.
		if err != nil {
			log.Fatalf("Fatal Server Shutdown: error with next_user_id in db (mostly overflow): %v", err)
		} else if id == 1 {
			if db.HGet(CTX, "user:1", "id").Val() == "1" {
				log.Fatalf("Fatal Server Shutdown: next_user_id in db is deleted.")
				return
			}
		}
		*/
		newuser.ID = strconv.FormatInt(id, 10)

		db.HSet(CTX, "users", newuser.Name, newuser.ID)

		hash, _ := bcrypt.GenerateFromPassword([]byte(newuser.Pwd), bcrypt.DefaultCost)
		newuser.Pwd = string(hash)

		db.GeoAdd(CTX, "locations", &redis.GeoLocation{
			Longitude: lon,
			Latitude:  lat,
			Name:      newuser.ID,
		})

		db.HSet(CTX, "user:"+newuser.ID, "id", newuser.ID, "name", newuser.Name, "dob", newuser.DOB, "address", newuser.Address, "latitude", lat, "longitude", lon, "description", newuser.Description, "createdAt", time.Now(), "pwd", newuser.Pwd) // use struct iterator or unroller.

		ctx.JSON(http.StatusCreated, &gin.H{"message": "account successfully created."})
	}
	return fx
}

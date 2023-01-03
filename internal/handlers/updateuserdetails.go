package handlers

import (
	"net/http"

	"github.com/Aadithya-V/user-system-microservice/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUserDetails(db *redis.Client) func(ctx *gin.Context) {
	fx := func(ctx *gin.Context) {
		/*
			Algorithm:
				Bind the received json to UpdatableUser struct.
				For-each field that has data,
					update the database with the corresponding field:value pair.
		*/
		// Binding data.
		var updatedUser models.UpdatableUser
		if err := ctx.BindJSON(&updatedUser); err != nil { // error not as I expected. Review. Extra fields ignored. Check if JSONpatch the only perfect soln?
			ctx.JSON(http.StatusBadRequest, &gin.H{"error": "Received information for unupdatable fields."}) // print json format for reference
			return
		}
		id := ctx.Param("id")
		// replace below unrolled loop with struct iterator.
		var tracker bool = false

		if updatedUser.Pwd != "" { // if pwd updated, delete all auth tokens, ie, log out all user's log-ins
			hash, _ := bcrypt.GenerateFromPassword([]byte(updatedUser.Pwd), bcrypt.DefaultCost)
			updatedUser.Pwd = string(hash)
			db.HSet(CTX, "user:"+id, "pwd", updatedUser.Pwd)
			tracker = true
		}
		if updatedUser.Description != "" {
			db.HSet(CTX, "user:"+id, "description", updatedUser.Description)
			tracker = true
		}
		if updatedUser.Address != "" {
			db.HSet(CTX, "user:"+id, "address", updatedUser.Address)
			tracker = true
		}
		if updatedUser.Latitude != 0.0 {
			db.HSet(CTX, "user:"+id, updatedUser.Latitude)
			db.GeoAdd(CTX, "locations", &redis.GeoLocation{
				Latitude: updatedUser.Latitude,
				Name:     id,
			})
			tracker = true
		}
		if updatedUser.Longitude != 0.0 {
			db.HSet(CTX, "user:"+id, updatedUser)
			db.GeoAdd(CTX, "locations", &redis.GeoLocation{
				Longitude: updatedUser.Longitude,
				Name:      id,
			})
			tracker = true
		}

		if tracker {
			ctx.JSON(http.StatusOK, &gin.H{"message": "data updated successfully."}) // need better error handling to report with 100% confidence
		} else {
			ctx.JSON(http.StatusBadRequest, &gin.H{"message:": "no data provided for updating."})
		}
	}

	return fx
}

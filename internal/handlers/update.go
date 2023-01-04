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
			ctx.JSON(http.StatusBadRequest, &gin.H{"error": "Incorrect updation information format."}) // print json format for reference
			return
		}
		id := ctx.Param("id")

		// replace below unrolled loop with struct iterator.
		tracker := false

		// Always check the coordinates first. If checked later after other db updates and coor is found to be invalid,
		// the entire updation process till then should have to performed as a transaction and rolled back.
		if updatedUser.Coordinates != NULLISLAND { // since float64 "comparable"
			// Validate coord
			lat, lon := updatedUser.Coordinates[0], updatedUser.Coordinates[1]

			if err := validateCoordinates(lat, lon); err != nil {
				ctx.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()}) // Received Invalid Coordinates
				return
			}

			db.GeoAdd(CTX, "locations", &redis.GeoLocation{
				Latitude:  lat,
				Longitude: lon,
				Name:      id,
			})
			// not adding coord to hash user:id coordinates as I might remove that filed soon.

			tracker = true
		}

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

		if tracker {
			ctx.JSON(http.StatusOK, &gin.H{"message": "data updated successfully."}) // need better error handling to report with 100% confidence
		} else {
			ctx.JSON(http.StatusBadRequest, &gin.H{"message:": "no data provided for updating."})
		}
	}

	return fx
}

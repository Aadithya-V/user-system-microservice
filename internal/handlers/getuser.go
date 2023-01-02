package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func GetUserByName(db *redis.Client) func(ctx *gin.Context) {
	fx := func(ctx *gin.Context) {
		// authentication
		token, err := ctx.Cookie("auth")
		if err == http.ErrNoCookie {
			ctx.JSON(http.StatusNotAcceptable, &gin.H{"err": "no auth cookie. Login and/or enable site cookies."})
			return //no cookie received
		}

		cid, err := db.HGet(CTX, "auths", token).Result() // id of the requester/client got from the auth key
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, &gin.H{"err": "wrong auth credential."})
			return //auth doesnt exist
		}

		tname := ctx.Param("name")                        // username of the target user to get info
		tid, err := db.HGet(CTX, "users", tname).Result() // id of the target user
		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, &gin.H{"err": "no such user exists."})
			return
		}

		var userdata *redis.SliceCmd

		if cid == tid {
			userdata = db.HMGet(CTX, "user:"+cid, "name", "description", "dob", "address", "latitude", "createdAt")
		} else {
			userdata = db.HMGet(CTX, "user:"+tid, "name", "description", "dob")
		}

		res, err := userdata.Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &gin.H{"error": "Internal server error in retreiving user data."})
			return
		}

		ctx.JSON(http.StatusOK, res)
	}

	return fx
}

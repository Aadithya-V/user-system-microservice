package handlers

// Simply delete entries of auth in auths and user::id and delete client cookie.

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func Logout(db *redis.Client) func(ctx *gin.Context) {
	fx := func(ctx *gin.Context) {
		token, err := ctx.Cookie("auth")
		if err == http.ErrNoCookie { //no cookie received
			ctx.JSON(http.StatusNotAcceptable, &gin.H{"err": "no auth cookie. Login and/or enable site cookies."})
			return
		}

		id, err := db.HGet(CTX, "auths", token).Result()
		if err != nil { //auth doesnt exist
			ctx.JSON(http.StatusUnauthorized, &gin.H{"err": "wrong auth credential."})
			return
		}
		// logout by deleting auth entries in db
		db.HDel(CTX, "auths", token)
		db.HDel(CTX, "user:"+id, "auth")
		// delete cookie
		ctx.SetCookie("auth", "", int(time.Unix(0, 0).Unix()), "", "", false, false)
		ctx.JSON(http.StatusAccepted, &gin.H{"message": "logged out"})
	}
	return fx
}

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
		if err == http.ErrNoCookie {
			return //no cookie received
		}

		res := db.HGet(CTX, "auths", token)
		if res.Err() != nil {
			return //auth doesnt exist
		}

		id := res.Val()
		// logout by deleting auth entries in db
		db.HDel(CTX, "auths", token)
		db.HDel(CTX, "user:"+id, "auth")
		// delete cookie
		ctx.SetCookie("auth", "", int(time.Unix(0, 0).Unix()), "", "", false, false)
		ctx.JSON(http.StatusAccepted, &gin.H{"logged out": "true"})
	}
	return fx
}

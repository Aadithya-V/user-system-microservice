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
		id := ctx.Param("id")
		token, _ := ctx.Cookie("auth") // no need to check error as tokenauth middleware has done it already.
		// logout by deleting current auth entry in db- TODO: del from set user:auths:id.
		db.HDel(CTX, "auths", token)
		db.HDel(CTX, "user:"+id, "auth")
		// delete cookie
		ctx.SetCookie("auth", "", int(time.Unix(0, 0).Unix()), "", "", false, false)
		ctx.JSON(http.StatusOK, &gin.H{"message": "logged out"})
	}
	return fx
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func TokenAuth(db *redis.Client) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("auth")
		if err == http.ErrNoCookie { //no cookie received
			ctx.JSON(http.StatusNotAcceptable, &gin.H{"err": "no auth cookie. Login and/or enable site cookies."})
			ctx.Abort()
			return
		}

		id, err := db.HGet(CTX, "auths", token).Result()

		if err == redis.Nil { //auth doesnt exist
			ctx.JSON(http.StatusUnauthorized, &gin.H{"err": "wrong auth credential."})
			ctx.Abort()
			return
		}

		ctx.AddParam("id", id)
		ctx.Next()
	}
}

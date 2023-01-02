package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Aadithya-V/user-system-microservice/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	tokenSize = 16 // size in Bytes
	CTX       = context.TODO()
)

type Logininfo struct {
	Uname string `json:"username" binding:"required"`
	Pwd   string `json:"pwd" binding:"required"`
}

func Login(db *redis.Client) func(ctx *gin.Context) {
	fx := func(ctx *gin.Context) {
		var login Logininfo
		if err := ctx.BindJSON(&login); err != nil {
			ctx.JSON(http.StatusInternalServerError, &gin.H{"error": "JSON Binding Failed"})
			return
		}
		id, err := db.HGet(CTX, "users", login.Uname).Result()
		if err != nil {
			ctx.JSON(http.StatusNotFound, &gin.H{"error": "Username Incorrect"})
			return
		}

		hash := db.HGet(CTX, "user:"+id, "pwd").Val()
		var token string
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(login.Pwd)); err == nil { // if already logged in update auth tokens
			token = auth.GenerateSecureToken(tokenSize)
			db.HSet(CTX, "user"+id, "auth", token)
			db.HSet(CTX, "auths", token, id)
		} else {
			ctx.JSON(http.StatusBadRequest, &gin.H{"error": "wrong password."})
			return
		}
		ctx.SetCookie("auth", token, int(time.Now().Unix()+3600*24*365), "", "", false, false)
		ctx.JSON(http.StatusOK, &gin.H{"message": "login successful."})

	}
	return fx
}

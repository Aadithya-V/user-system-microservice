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
			//ctx.JSON(http., struct{}{})
			return
		}
		res := db.HGet(CTX, "users", login.Uname)
		if res.Err() != nil {
			return
		}
		id := res.Val()

		hash := db.HGet(CTX, "user:"+id, "pwd").Val()
		var token string
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(login.Pwd)); err == nil { // && not already logged in, ie, no auth token exists
			token = auth.GenerateSecureToken(tokenSize)
			db.HSet(CTX, "user"+id, "auth", token)
			db.HSet(CTX, "auths", token, id)
		}
		ctx.SetCookie("auth", token, int(time.Now().Unix()+3600*24*365), "", "", false, false)
		ctx.JSON(http.StatusOK, struct{}{}) //add to header/cookie

	}
	return fx
}

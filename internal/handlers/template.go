package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func name(db *redis.Client) func(ctx *gin.Context) {
	fx := func(ctx *gin.Context) {

	}

	return fx
}

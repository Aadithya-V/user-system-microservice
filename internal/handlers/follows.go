package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func FollowUser(db *redis.Client) func(ctx *gin.Context) {
	fx := func(ctx *gin.Context) {
		nametoFollow := ctx.Param("name")
		idtoFollow, err := db.HGet(CTX, "users", nametoFollow).Result()
		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, &gin.H{"error": "Username incorrect. Does not exist."})
		}
		cid := ctx.Param("id") // id of user following idtoFollow

		if cid == idtoFollow {
			ctx.JSON(http.StatusForbidden, &gin.H{"error": "You are not allowed to follow/unfollow yourself."})
			return
		}

		db.ZAdd(CTX, "followers:"+idtoFollow, redis.Z{Score: float64(time.Now().Unix()), Member: cid})
		db.ZAdd(CTX, "following:"+cid, redis.Z{Score: float64(time.Now().Unix()), Member: idtoFollow})

		ctx.JSON(http.StatusOK, &gin.H{"message": "..followed"})
	}

	return fx
}

// Logical inverse of func FollowUser(). Same code repeated instead of handling both follow and unfollow requests
// in one function itself (by checking a bool value passed as parameter) as -
// the conditional check is too random for CPU's speculative execution
// resulting in too high a penalty for frequent misses.
// Wait, :( statistically follows are more frequent than unfollows...
// But still, will take the minuscle performance improvement by avoiding the penalty altogether..
func UnfollowUser(db *redis.Client) func(ctx *gin.Context) {
	fx := func(ctx *gin.Context) {
		nametoUnfollow := ctx.Param("name")
		idtoUnfollow, err := db.HGet(CTX, "users", nametoUnfollow).Result()
		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, &gin.H{"error": "Username incorrect. Does not exist."})
		}
		cid := ctx.Param("id") // id of user unfollowing idtoUnfollow

		if cid == idtoUnfollow {
			ctx.JSON(http.StatusForbidden, &gin.H{"error": "You are not allowed to follow/unfollow yourself."})
			return
		}

		db.ZRem(CTX, "followers:"+idtoUnfollow, cid)
		db.ZRem(CTX, "following:"+cid, idtoUnfollow)

		ctx.JSON(http.StatusOK, &gin.H{"message": "..Unfollowed"})
	}

	return fx
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Aadithya-V/user-system-microservice/internal/database"
	"github.com/Aadithya-V/user-system-microservice/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

// Constants specifying the listening addresses of
// the redis server and the gin router engine.
var (
	ListenAddr = "localhost:8080"
	RedisAddr  = "localhost:6379"
)

func main() {
	// Initialize database connection
	db, err := database.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "", //get pwd from env
		DB:       0,
	})
	if err != nil {
		log.Fatalf("Database Connection Failed: %v", err)
	}

	// Initialize router gin engine
	router := initRouter(db)

	//
	startHttpServer(router) // blocking. Spin this off into a go routine if subsequent code is added.

}

// Function initRouter() initialises a router,
// maps the routes and returns a pointer to it
// which is *gin.Engine
func initRouter(db *redis.Client) *gin.Engine {
	router := gin.Default()

	// Routes mapping
	// add logout, getuser, update, follow, unfollow,followers, following, nearby to authorization group.
	authorized := router.Group("/user", handlers.TokenAuth(db))
	authorized.POST("/logout", handlers.Logout(db))
	authorized.GET("/:name", handlers.GetUserByName(db))
	authorized.PUT("/update", handlers.UpdateUserDetails(db))
	authorized.POST("/:name/follow", handlers.FollowUser(db))
	authorized.POST("/:name/unfollow", handlers.UnfollowUser(db))

	router.POST("/register", handlers.Register(db))
	router.POST("/login", handlers.Login(db))

	//router.GET("/user/:name/followers", )	// GET("/user/:id") too should retrieve follows.
	//router.GET("/user/:name/following", )
	//router.GET("/user/nearby")
	return router
}

func startHttpServer(router *gin.Engine) {
	srv := &http.Server{
		Addr:    ListenAddr,
		Handler: router,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		// for graceful shutdown..
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listener(s), or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	// blocking service of connections
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

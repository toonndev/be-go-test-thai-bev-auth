package router

import (
	"os"

	"be-go-test-thai-bev-auth/internal/handler"
	"be-go-test-thai-bev-auth/internal/middleware"
	"be-go-test-thai-bev-auth/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(authHandler *handler.AuthHandler, blacklistRepo repository.TokenBlacklistRepository) *gin.Engine {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	jwtMiddleware := middleware.JWTMiddleware(blacklistRepo)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", jwtMiddleware, authHandler.Me)
			auth.POST("/logout", jwtMiddleware, authHandler.Logout)
		}
	}

	return r
}

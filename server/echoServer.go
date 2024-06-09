package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sudhir512kj/ecommerce_backend/config"
	"github.com/sudhir512kj/ecommerce_backend/database"
	"github.com/sudhir512kj/ecommerce_backend/internal/handlers"
	"github.com/sudhir512kj/ecommerce_backend/internal/repository"
)

type echoServer struct {
	app  *gin.Engine
	db   database.Database
	conf *config.Config
}

func NewEchoServer(conf *config.Config, db database.Database) Server {
	ginApp := gin.Default()
	// echoApp.Logger.SetLevel(log.DEBUG)

	return &echoServer{
		app:  ginApp,
		db:   db,
		conf: conf,
	}
}

func (s *echoServer) Start() {
	// s.app.Use(middleware.Recover())
	// s.app.Use(middleware.Logger())

	// Create the repositories
	userRepo := repository.NewUserRepository(s.db.GetDb())

	// Create the handlers
	userHandler := handlers.NewUserHandler(userRepo)

	// Define the API routes
	api := s.app.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.POST("/forgot-password", userHandler.ForgotPassword)
			users.POST("/reset-password", userHandler.AuthMiddleware, userHandler.ChangePassword)
			// users.GET("/profile", userHandler.GetProfile)
			users.PUT("/update-profile", userHandler.AuthMiddleware, userHandler.UpdateProfile)
			users.POST("/verify-otp", userHandler.VerifyOTP)
		}
	}

	// Health check adding
	s.app.GET("/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Run(serverUrl)
}

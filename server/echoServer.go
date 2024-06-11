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
	productRepo := repository.NewProductRepository(s.db.GetDb())
	categoryRepo := repository.NewCategoryRepository(s.db.GetDb())
	productImageRepo := repository.NewProductImageRepository(s.db.GetDb())
	cartRepo := repository.NewCartRepository(s.db.GetDb())
	orderRepo := repository.NewOrderRepository(s.db.GetDb())
	paymentRepo := repository.NewPaymentRepository(s.db.GetDb())

	// Create the handlers
	userHandler := handlers.NewUserHandler(userRepo)
	productHandler := handlers.NewProductHandler(productRepo, categoryRepo, productImageRepo)
	cartHandler := handlers.NewCartHandler(cartRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo, cartRepo, paymentRepo)
	// paymentHandler := handlers.NewPaymentHandler(paymentRepo)

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
		products := api.Group("/products")
		{
			products.GET("/", productHandler.GetAllProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.POST("/", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}
		carts := api.Group("/cart")
		{
			carts.POST("/", userHandler.AuthMiddleware, cartHandler.AddToCart)
			carts.GET("/", userHandler.AuthMiddleware, cartHandler.GetCartInfo)
			carts.PUT("/:id", cartHandler.UpdateCart)
			carts.DELETE("/:id", cartHandler.DeleteFromCart)
			carts.PUT("/:id/save-for-later", cartHandler.SaveForLater)
		}
		orders := api.Group("/orders")
		{
			// add endpoints for all handlers defined in order.go in handlers
			orders.POST("/", orderHandler.Checkout)
		}
	}

	// Health check adding
	s.app.GET("/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Run(serverUrl)
}

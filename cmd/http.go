package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"github.com/mhasnanr/ewallet-ums/external"
	"github.com/mhasnanr/ewallet-ums/helpers"
	handler "github.com/mhasnanr/ewallet-ums/internal/handler/http"
	"github.com/mhasnanr/ewallet-ums/internal/middleware"
	"github.com/mhasnanr/ewallet-ums/internal/repository"
	"github.com/mhasnanr/ewallet-ums/internal/services"
	"gorm.io/gorm"
)

func ServeHTTP(db *gorm.DB) {
	r := gin.New()

	r.Use(middleware.LoggerMiddleware(bootstrap.Log))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "server is healthy"})
	})

	walletGRPCClient, grpcConn, err := external.NewWalletGRPC()
	if err != nil {
		log.Fatalf("failed to initialized gRPC: %v", err)
	}

	defer grpcConn.Close()

	jwtManager := &helpers.JWTManager{}
	passwordHasher := &helpers.PasswordHasher{}

	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db)

	authMiddleware := middleware.NewAuthMiddleware(sessionRepository, jwtManager)
	userService := services.NewUserService(userRepository, sessionRepository, jwtManager, passwordHasher, walletGRPCClient)
	userHandler := handler.NewUserHandler(userService, authMiddleware)

	userHandler.RegisterRoute(r)

	httpPort := bootstrap.GetEnv("HTTP_PORT", "8080")
	server := &http.Server{Addr: ":" + httpPort, Handler: r}

	fmt.Printf("http server is running on port %s...\n", httpPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server stopped")
	}
}

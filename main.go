package main

import (
	"auth-gorm-echo/auth"
	"auth-gorm-echo/campaign"
	"auth-gorm-echo/config"
	"auth-gorm-echo/handler"
	"auth-gorm-echo/helper"
	"auth-gorm-echo/user"
	"net/http"
	"strings"

	// "fmt"
	// "net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Custom Error
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Connect to database
	config.DatabaseInit()
	db := config.GetDB()

	dbGorm, err := db.DB()
	if err != nil {
		panic(err)
	}
	dbGorm.Ping()

	// initialize redis
	config.RedisInit()

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	
	router := echo.New()
	router.Validator = &CustomValidator{validator: validator.New()}

	// access images
	router.Static("/images", "./images")

	// Router
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)

	api.GET("/campaigns", campaignHandler.GetCampaigns)

	api.Use(authMiddleware(authService, userService))
	api.GET("/users/fetch", userHandler.FetchUser)
	api.POST("/avatars", userHandler.UploadAvatar)


	router.Logger.Fatal(router.Start(":9000"))
}

// middleware
func authMiddleware(authService auth.Service, userService user.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if !strings.Contains(authHeader, "Bearer") {
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				return c.JSON(http.StatusUnauthorized, response)
			}

			tokenString := ""
			arrayToken := strings.Split(authHeader, " ")
			if len(arrayToken) == 2 {
				tokenString = arrayToken[1]
			}

			token, err := authService.ValidateToken(tokenString)
			if err != nil {
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				return c.JSON(http.StatusUnauthorized, response)
			}

			claim, ok := token.Claims.(jwt.MapClaims)

			if !ok || !token.Valid {
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				return c.JSON(http.StatusUnauthorized, response)
			}

			userID := int(claim["user_id"].(float64))

			user, err := userService.GetUserByID(userID)
			if err != nil {
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				return c.JSON(http.StatusUnauthorized, response)
			}

			c.Set("currentUser", user)
			return next(c)
		}
	}
}

// input dari user
// handler mapping input dari user ke struct input
// service mapping ke struct User
// repository save struct User ke db
// db
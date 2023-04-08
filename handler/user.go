package handler

import (
	"auth-gorm-echo/auth"
	"auth-gorm-echo/config"
	"auth-gorm-echo/helper"
	"auth-gorm-echo/user"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

type RequestRedis struct {
	ID string
	Token string
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c echo.Context) error {
	// tangkap input dari user
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter service

	var input user.RegisterUserInput
	err := c.Bind(&input)
	if err := c.Validate(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := echo.Map{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	return c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c echo.Context) error {
	// user memasukan input (email & password)
	// input ditangkap handler
	// mapping dari input user ke input struct
	// input struct passing service
	// di service mencari dg bantuan repository user dengan email x
	// mencocokan password

	var input user.LoginInput

	err := c.Bind(&input)
	if err := c.Validate(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := echo.Map{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	loggedInUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := echo.Map{"errors": err.Error()}

		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", errorMessage)
		return c.JSON(http.StatusBadRequest, response)
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Redis Session
	sessionExp := config.GetsessionExp()
	redisCtx := config.GetRedisCtx()
	rdb := config.RedisConnect()

	userID := strconv.Itoa(loggedInUser.ID)

	reqRedis := RequestRedis{
		ID: userID,
		Token: token,
	}
	req, _ := json.Marshal(reqRedis)
	
	sessionID := fmt.Sprintf("session:%d", loggedInUser.ID)
	err = rdb.Set(redisCtx, sessionID, req, sessionExp).Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error saving session"})
	}

	formatter := user.FormatUser(loggedInUser, token)

	response := helper.APIResponse("Successfuly logged in", http.StatusOK, "success", formatter)

	return c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c echo.Context) error {
	// input email dari user
	// input email di mapping ke struct input
	// struct input di passing ke service
	// service akan memanggil repository - email sudah ada atau belum
	// repository - db

	var input user.CheckEmailInput
	
	err := c.Bind(&input)
	if err := c.Validate(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := echo.Map{"errors": errors}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := echo.Map{"errors": "Server error"}

		response := helper.APIResponse("Email checking failed", http.StatusBadRequest, "error", errorMessage)
		return c.JSON(http.StatusBadRequest, response)
	}

	data := echo.Map{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	return c.JSON(http.StatusOK, response)
}

func (h *userHandler) FetchUser(c echo.Context) error {
	currentUser := c.Get("currentUser").(user.User)

	formatter := user.FormatUser(currentUser, "")

	response := helper.APIResponse("Successfuly fetch user data", http.StatusOK, "success", formatter)

	return c.JSON(http.StatusOK, response)
}
package handler

import (
	"auth-gorm-echo/campaign"
	"auth-gorm-echo/helper"
	"auth-gorm-echo/user"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// tangkap parameter di handler
// handler ke service
// service menentukan repository mana yg dipanggil
// repository: FindAll, FindByUserId
// db

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c echo.Context) error {

	userID, _ := strconv.Atoi(c.QueryParam("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	return c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c echo.Context) error {
	// api/v1/campaigns/:id
	// handler : mapping id dari url ke struct input -> service, call formatter
	// service : inputnya struct input -> menangkap id di url, manggil repository -> call formatter
	// repository : get campaign by id

	var input campaign.GetCampaignDetailInput

	err := c.Bind(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	return c.JSON(http.StatusOK, response)
}

// tangkap parameter dari user ke input struct
// ambil current user dari jwt/handler
// panggil service, paramternya input struck (dan juga buat slug)
// panggil repository untuk simpan ke db

func (h *campaignHandler) CreateCampaign(c echo.Context) error {
	var input campaign.CreateCampaignInput

	err := c.Bind(&input)
	if err := c.Validate(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := echo.Map{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	currentUser := c.Get("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	response := helper.APIResponse("Campaign has been created", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	return c.JSON(http.StatusOK, response)
}
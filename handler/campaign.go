package handler

import (
	"auth-gorm-echo/campaign"
	"auth-gorm-echo/helper"
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
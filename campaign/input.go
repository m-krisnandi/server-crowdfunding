package campaign

type GetCampaignDetailInput struct {
	ID int `param:"id" validate:"required"`
}
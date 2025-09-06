package dto

type ChatDto struct {
	Query   string `json:"message" form:"message" binding:"required"`
	Remarks string `json:"remarks" form:"remarks"`
}

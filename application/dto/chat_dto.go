package dto

type ChatDto struct {
	Query   string `json:"query" form:"query"`
	Remarks string `json:"remarks" form:"remarks"`
}

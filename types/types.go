package types

type Event struct{
	Type string `json:"type"`
	CreatedAt string `json:"created_at"`
	PageUrl string `json:"page_url"`
	UserId string `json:"user"`
}

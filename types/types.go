package types

import "time"

type Event struct{
	Type string `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	PageUrl string `json:"page_url"`
	User struct{
		Id string
	}`json:"user"`
}

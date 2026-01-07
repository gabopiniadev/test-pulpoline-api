package models

import "context"

type Request struct {
	ID      string
	Text    string
	Context context.Context
	Result  chan Response
	Error   chan error
}

type Response struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Response string `json:"response"`
}

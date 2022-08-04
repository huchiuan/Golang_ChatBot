package model

import "time"

type Config struct {
	LineChannelSecret string
	LineChannelToken  string
	Port              string
	MongoAccount      string
	MongoPassword     string
}

type LineMessage struct {
	UserID  string    `json:"userid"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type ApiRequestBody struct {
	Message string
	UserID  string
}

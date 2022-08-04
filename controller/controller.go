package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	configpackage "golang_chatbot/config"
	model "golang_chatbot/model"
	sqlbublic "golang_chatbot/sqlpublic"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func QueryUserMessages(c *gin.Context) {
	sqlbublic.ConnectToDB()
	var requestbody model.ApiRequestBody

	if err := c.BindJSON(&requestbody); err != nil {
		log.Fatal(err)
	}

	fmt.Println(requestbody.UserID)

	messages := sqlbublic.GetMessages(requestbody.UserID)

	c.JSON(200, messages)
}

func PushMessage(c *gin.Context) {

	var requestbody model.ApiRequestBody

	if err := c.BindJSON(&requestbody); err != nil {
		log.Fatal(err)
	}

	fmt.Println(requestbody.Message)

	config, _ := configpackage.InitConfig()
	bot, err := linebot.New(config.LineChannelSecret, config.LineChannelToken)
	_, err = bot.BroadcastMessage(linebot.NewTextMessage(requestbody.Message)).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func ReceiveMessage(c *gin.Context) {
	sqlbublic.ConnectToDB()
	var err error
	config, err := configpackage.InitConfig()

	bot, err := linebot.New(config.LineChannelSecret, config.LineChannelToken)
	events, err := bot.ParseRequest(c.Request)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		UserID := event.Source.UserID
		Time := event.Timestamp
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			case *linebot.TextMessage:

				sqlbublic.SaveMessage(UserID, message.Text, Time)
				TimeToString := event.Timestamp.String()
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("msg ID:"+message.ID+"\n\n文字內容:"+message.Text+"\n\n使用者ID: "+UserID+"\n\n時間: "+TimeToString)).Do(); err != nil {
					log.Print(err)
				}

			}
		}
	}
}

package controller

import (
	"github.com/gin-gonic/gin"

	configpackage "golang_chatbot/config"
	model "golang_chatbot/model"
	sqlbublic "golang_chatbot/sqlpublic"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func ReceiveMessage(c *gin.Context) {
	client := sqlbublic.ConnectToDB()

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
		userID := event.Source.UserID
		time := event.Timestamp
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			case *linebot.TextMessage:

				sqlbublic.SaveMessage(client, userID, message.Text, time)
				timeToString := event.Timestamp.String()
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("msg ID:"+message.ID+"\n\n文字內容:"+message.Text+"\n\n使用者ID: "+userID+"\n\n時間: "+timeToString)).Do(); err != nil {
					log.Print(err)
				}

			}
		}
	}
}

func PushMessage(c *gin.Context) {

	var requestBody model.ApiRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		log.Fatal(err)
	}

	config, _ := configpackage.InitConfig()
	bot, err := linebot.New(config.LineChannelSecret, config.LineChannelToken)
	_, err = bot.BroadcastMessage(linebot.NewTextMessage(requestBody.Message)).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func QueryUserMessages(c *gin.Context) {
	client := sqlbublic.ConnectToDB()
	var requestBody model.ApiRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		log.Fatal(err)
	}

	messages := sqlbublic.GetMessages(client, requestBody.UserID)

	c.JSON(200, messages)
}

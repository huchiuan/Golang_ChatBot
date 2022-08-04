package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	configpackage "golang_chatbot/config"
	sqlbublic "golang_chatbot/sqlpublic"
	"log"
	"net/http"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
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
			// Handle only on text message
			case *linebot.TextMessage:
				// GetMessageQuota: Get how many remain free tier push message quota you still have this month. (maximum 500)
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				// message.ID: Msg unique ID
				// message.Text: Msg text

				sqlbublic.SaveMessage(UserID, message.Text, Time)
				TimeToString := event.Timestamp.String()
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("msg ID:"+message.ID+":"+"Get:"+message.Text+"使用者ID: "+UserID+"時間: "+TimeToString+" , \n OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
					log.Print(err)
				}

			// Handle only on Sticker message
			case *linebot.StickerMessage:
				var kw string
				for _, k := range message.Keywords {
					kw = kw + "," + k
				}

				outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s  text: %s ", message.StickerID, message.PackageID, kw, message.Text)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

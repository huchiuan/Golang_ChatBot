package config

import (
	model "golang_chatbot/model"

	"github.com/spf13/viper"
)

func InitConfig() (config model.Config, err error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err = viper.ReadInConfig()
	if err != nil {
		panic("讀取設定檔出現錯誤，原因為：" + err.Error())
	}

	config = model.Config{
		LineChannelSecret: viper.GetString("Line.ChannelSecret"),
		LineChannelToken:  viper.GetString("Line.ChannelToken"),
		MongoAccount:      viper.GetString("MongoDB.MongoAccount"),
		MongoPassword:     viper.GetString("MongoDB.MongoPassword"),
	}

	return config, err
}

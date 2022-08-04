// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	configpackage "golang_chatbot/config"
	controller "golang_chatbot/controller"

	// sqlbublic "golang_chatbot/sqlpublic"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func main() {

	router := gin.Default()

	router.POST("/callback", controller.ReceiveMessage)
	router.GET("/", controller.IndexHandler)

	config, _ := configpackage.InitConfig()
	port := config.Port
	addr := fmt.Sprintf(":%s", port)

	router.Run(addr)

}

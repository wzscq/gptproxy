package main

import (
	"gptproxy/common"
	"gptproxy/mqtt"
	"gptproxy/customer"
	"log"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"os"
)

func main() {
	//初始化时区
    var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	//设置log打印文件名和行号
    log.SetFlags(log.Lshortfile | log.LstdFlags)

	confFile:="conf/conf.json"
    if len(os.Args)>1 {
        confFile=os.Args[1]
        log.Println(confFile)
    }

    //初始化配置
    conf:=common.InitConfig(confFile)

	//初始化消息处理器
	messageHandler:=&customer.MessageHandler{}
	//mqttclient
	mqttClient:=mqtt.MQTTClient{
		Broker:conf.MQTT.Broker,
		User:conf.MQTT.User,
		Password:conf.MQTT.Password,
		BillTopic:conf.MQTT.BillTopic,
		RedirectTopic:conf.MQTT.RedirectTopic,
		Handler:messageHandler,
		ClientID:conf.MQTT.ClientID,
	}
	mqttClient.Init()
	
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowAllOrigins:true,
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"*"},
        AllowCredentials: true,
    }))

	router.Run(conf.Service.Port)
}
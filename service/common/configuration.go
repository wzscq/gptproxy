package common

import (
	"log"
	"os"
	"encoding/json"
)

type serviceConf struct {
	Port string `json:"port"`
}

type mysqlConf struct {
	Server string `json:"server"`
	Password string `json:"password"`
	User string `json:"user"`
	DBName string `json:"dbName"`
	ConnMaxLifetime int `json:"connMaxLifetime"` 
  MaxOpenConns int `json:"maxOpenConns"`
  MaxIdleConns int `json:"maxIdleConns"`
}

type AccountCacheConf struct {
	Server string `json:"server"`
	Password string `json:"password"`
	DB int `json:"db"`
}

type mqttConf struct {
	Broker string `json:"broker"`
	User string `json:"user"`
	Password string `json:"password"`
	ClientID string `json:"clientID"`
	RedirectTopic string `json:"redirectTopic"`
	BillTopic string `json:"billTopic"`
}

type Config struct {
	Service serviceConf `json:"service"`
	AccountCache AccountCacheConf `json:"accountCache"`
	MQTT mqttConf `json:"mqtt"`
	Mysql  mysqlConf  `json:"mysql"`
}

var gConfig Config

func InitConfig(confFile string)(*Config){
	log.Println("init configuation start ...")
	//获取用户账号
	//获取用户角色信息
	//根据角色过滤出功能列表
	fileName := confFile
	filePtr, err := os.Open(fileName)
	if err != nil {
        log.Fatal("Open file failed [Err:%s]", err.Error())
    }
    defer filePtr.Close()

	// 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&gConfig)
	if err != nil {
		log.Println("json file decode failed [Err:%s]", err.Error())
	}
	log.Println("init configuation end")
	return &gConfig
}

func GetConfig()(*Config){
	return &gConfig
}
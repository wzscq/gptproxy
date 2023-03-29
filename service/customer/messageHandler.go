package customer

import (
	"encoding/json"
	"log"
	"strings"
)

type WeiChatMessage struct {
	ToUserName string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	CreateTime string `json:"CreateTime"`
	MsgType string `json:"MsgType"`
	Event string `json:"Event"`
	EventKey string `json:"EventKey"`
	Content string `json:"Content"`
	MsgId string `json:"MsgId"`
	MsgDataId string `json:"MsgDataId"`
	Idx string `json:"Idx"`
}

type BillRec struct {
	Sessionid string `json:"sessionid"`
	Amount int `json:"amount"`
}

type MqttClient interface {
	Publish(topic,content string)(int)
}

type MessageHandler struct {
	Repo Repository
	AccountCache *AccountCache
	MqttClient MqttClient
}

func (m *MessageHandler)DealWeiChatMessage(msg []byte){
	var message WeiChatMessage
	err:=json.Unmarshal(msg,&message)
	if err!=nil {
		log.Println(err)
		return
	}

	log.Printf("message type:%s\n",message.MsgType)
	
	//更新用户状态
	if message.MsgType=="event" {
		log.Printf("event type:%s\n",message.Event)
		if message.Event=="subscribe" {
			err:=m.Repo.AddCustomer(message.FromUserName)
			if err != nil {
				log.Println(err)
			}
			//如果是关注账号的话，需要更新缓存
			total,useed,err:=m.Repo.GetUserToken(message.FromUserName)
			if err != nil {
				log.Println(err)
			} else {
				m.AccountCache.AddAccount(message.FromUserName,total,useed)
			}
		}else if message.Event=="unsubscribe" {
			err:=m.Repo.DeactiveCustomer(message.FromUserName)
			if err != nil {
				log.Println(err)
			}
			//如果是取消关注账号的话，需要从缓存中删除
			m.AccountCache.RemoveAccount(message.FromUserName)
		}

		//如果给了扫码场景信息，那么就发送一个消息到MQ
		if message.EventKey!="" {
			eventKey:=message.EventKey
			if strings.Contains(eventKey,"qrscene_") {
				eventKey=strings.Replace(eventKey,"qrscene_","",-1)
			}
			m.MqttClient.Publish("qrlogin/"+eventKey,message.FromUserName)
		}

		return
	}
}

func (m *MessageHandler)DealBillRecMessage(msg []byte){
	var billRec BillRec
	err:=json.Unmarshal(msg,&billRec)
	if err!=nil {
		log.Println(err)
		return
	}
	
	//更新用户状态
	err=m.Repo.IncreaseUsedToken(billRec.Sessionid,billRec.Amount)
	if err != nil {
		log.Println(err)
		return
	}

	//更新缓存
	m.AccountCache.IncreaseUsedToken(billRec.Sessionid,billRec.Amount)
}
package customer

import (
	"encoding/json"
	"log"
)

type WeiChatMessage struct {
	ToUserName string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	CreateTime string `json:"CreateTime"`
	MsgType string `json:"MsgType"`
	Event string `json:"Event"`
	Content string `json:"Content"`
	MsgId string `json:"MsgId"`
	MsgDataId string `json:"MsgDataId"`
	Idx string `json:"Idx"`
}

type BillRec struct {
	Sessionid string `json:"sessionid"`
	Amount int `json:"amount"`
}

type MessageHandler struct {
}

func (m *MessageHandler)DealWeiChatMessage(msg []byte){
	var message WeiChatMessage
	err:=json.Unmarshal(msg,&message)
	if err!=nil {
		log.Println(err)
		return
	}
	
	//更新用户状态
}

func (m *MessageHandler)DealBillRecMessage(msg []byte){
	var billRec BillRec
	err:=json.Unmarshal(msg,&billRec)
	if err!=nil {
		log.Println(err)
		return
	}
	log.Println(billRec)
}
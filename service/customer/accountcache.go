package customer

import (
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

type AccountCache struct {
	Server string
	Password string 
	DB int 
}

func (this *AccountCache)GetRedisClient()(*redis.Client){
	//初始化redis
	return redis.NewClient(&redis.Options{
		Addr:     this.Server,
		Password: this.Password, 
		DB:       this.DB,  
	})
}

func (this *AccountCache)IncreaseUsedToken(key string,token int){
	client:=this.GetRedisClient()
	defer client.Close()

	//获取当前已使用的token
	usedTokenStr, err:=client.Get(client.Context(), key+":usedToken").Result()
	if err != nil {
		log.Println(err)
		return
	}

	//将totalToken转换为int类型
	usedToken,err:= strconv.Atoi(usedTokenStr)
	if err != nil {
		log.Println(err)
		return
	}

	//增加已使用的token
	usedToken+=token

	//设置已使用的token
	err=client.Set(client.Context(), key+":usedToken", usedToken, 0).Err()
	if err != nil {
		log.Println(err)
		return
	}
}

func (this *AccountCache)RemoveAccount(key string){
	client:=this.GetRedisClient()
	defer client.Close()

	//删除总token
	err:=client.Del(client.Context(), key+":totalToken").Err()
	if err != nil {
		log.Println(err)
		return
	}

	//删除已使用的token
	err=client.Del(client.Context(), key+":usedToken").Err()
	if err != nil {
		log.Println(err)
		return
	}
}

func (this *AccountCache)AddAccount(key string,token int,usedToken int){
	client:=this.GetRedisClient()
	defer client.Close()

	//设置总token
	err:=client.Set(client.Context(), key+":totalToken", token, 0).Err()
	if err != nil {
		log.Println(err)
		return
	}

	//设置已使用的token
	err=client.Set(client.Context(), key+":usedToken", usedToken, 0).Err()
	if err != nil {
		log.Println(err)
		return
	}
}

func (this *AccountCache)GetToken(key string)(int){
	client:=this.GetRedisClient()
	defer client.Close()

	totalTokenStr, err:=client.Get(client.Context(), key+":totalToken").Result()
	if err != nil {
		log.Println(err)
		return 0
	}

	//将totalToken转换为int类型
	totalToken,err:= strconv.Atoi(totalTokenStr)
	if err != nil {
		log.Println(err)
		return 0
	}

	usedTokenStr, err:=client.Get(client.Context(), key+":usedToken").Result()
	if err != nil {
		log.Println(err)
		return 0
	}

	//将totalToken转换为int类型
	usedToken,err:= strconv.Atoi(usedTokenStr)
	if err != nil {
		log.Println(err)
		return 0
	}

	return totalToken - usedToken
}
package main

import (
	"fmt"
	"net/http"
	"wechat/context"
	"wechat/server"
)

const (
	Token          = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	CorpId         = "xxxxxxxxxxxxxxxxxxxxxxxxx"
	AppSecret      = "xxxxxxxxxxxxxxxxxxxxxxxxxxx"
	EncodingAESKey = "xxxxxxxxxxxxxxxxxxxxxxxxxxxx"
)

func wechat(rw http.ResponseWriter, req *http.Request) {
	context := &context.Context{
		Request:        req,
		Writer:         rw,
		AppID:          CorpId,
		AppSecret:      AppSecret,
		Token:          Token,
		EncodingAESKey: EncodingAESKey,
	}
	err := ser.Serve()
	if err != nil {
		log.Printf("Handle the message error: ", err)
	}
}

func main() {
	http.HandleFunc("/wechat", wechat)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}

}

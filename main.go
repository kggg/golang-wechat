package main

import (
	"fmt"
	"github.com/kggg/golang-wechat/utils/msgcrypt"
	"log"
	"net/http"
)

const (
	Token          = "got from wechat access token"
	CorpId         = "company wechat id"
	EncodingAesKey = "jxD35KsNDtDBWUGgr4rs0dddddddddddddddddddddd"
)

func wechat(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	MsgSig := req.Form.Get("msg_signature")
	TimeStamp := req.Form.Get("timestamp")
	Nonce := req.Form.Get("nonce")
	EchoStr := req.Form.Get("echostr")

	valid := msgcrypt.ValidateMsg(Token, TimeStamp, Nonce, EchoStr, MsgSig)
	if valid {
		fmt.Println("validate ok")
		xml_content, err := msgcrypt.DecryptMsg(EchoStr, EncodingAesKey)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Fprintf(rw, string(xml_content))
	} else {
		fmt.Fprintf(rw, "valid fail")
	}

	if req.Method == "POST" {

	}

}

func main() {
	http.HandleFunc("/", wechat)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}

}

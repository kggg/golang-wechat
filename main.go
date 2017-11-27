package main

import (
	"fmt"
	"log"
	"net/http"
	"wechat/utils/msgcrypt"
)

const (
	Token          = "Tj7rlp8BSCS7SFEn8rE9CSep5Dbl2lmx"
	CorpId         = "ww3bedc4cae75db7ec"
	EncodingAesKey = "jxD35KsNDtDBWUGgr4rs0JPEdYO4KNU7CmcBzdTBugK"
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

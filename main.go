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
		body := req.Body
		xml_content, err := msgcrypt.DecryptMsg(req.Encrypt, EncodingAesKey)
		if err != nil {
			log.Fatalln(err)
		}
		var rec MsgBody
		err = xml.Unmarshal(xml_content, &rec)
		if err != nil {
			beego.Debug("msgbody invalid:", err)
			log.Fatalln(err)
		}
		var sendmsg CDATA
		if rec.Content == "你好" {
			sendmsg = "我很好"
		} else if rec.Content == "在忙什么" {
			sendmsg = "处理国家大事，国家不能没有我"
		} else if rec.Content == "在干嘛" {
			sendmsg = "想看看美女， 但好多国家大事要跟进"
		} else if rec.Content == "在吗" {
			sendmsg = "我说我在还是不在呢？"
		} else {
			sendmsg = "我是聊天机器人，我有什么可以帮你吗"
		}
		msg := &MsgBody{
			ToUserName:   rec.FromUserName,
			FromUserName: rec.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      rec.MsgType,
			//Content:      "hello world",
			Content: sendmsg,
		}
		msg_xml, err := xml.MarshalIndent(msg, "", "  ")
		if err != nil {
			beego.Debug("xml marshal msg error: ", err)
		}
		ciphertext, err := msgcrypt.EncryptMsg(msg_xml, EncodingAesKey, CorpId)
		if err != nil {
			beego.Debug("xml encrypt msg error: ", err)
		}
		nowtime := time.Now().Unix()
		stime := strconv.FormatInt(nowtime, 10)
		msg_Sig := msgcrypt.MakeSHA1Slice(Token, stime, Nonce, ciphertext)
		resp := &ResponseBody{
			Nonce:        CDATATEXT{Nonce},
			Encrypt:      CDATATEXT{ciphertext},
			TimeStamp:    nowtime,
			MsgSignature: CDATATEXT{msg_Sig},
		}
		output, err := xml.MarshalIndent(resp, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Fprintf(c.Ctx.ResponseWriter, string(output))

	}

}

func main() {
	http.HandleFunc("/", wechat)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}

}

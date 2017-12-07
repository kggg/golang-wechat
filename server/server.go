package server

import (
	"fmt"
	"log"
	"wechat/context"
	"wechat/utils/msgcrypt"
)

//Server struct
type Server struct {
	*context.Context

	requestRawXMLMsg  []byte
	requestMsg        []byte
	responseRawXMLMsg []byte
	responseMsg       interface{}

	random    []byte
	nonce     string
	timestamp int64
}

//NewServer init
func NewServer(context *context.Context) *Server {
	srv := new(Server)
	srv.Context = context
	return srv
}

//Serve 处理微信的请求消息
func (srv *Server) Serve() {

	if srv.Request.Method == "GET" {
		echostr, exists := srv.QueryURLParam("echostr")
		if !exists {
			fmt.Errorf("echostr not exists in url")
		}
		if !srv.Validate(echostr) {
			fmt.Errorf("消息不合法，验证签名失败")
		}
		xml_content, err := msgcrypt.DecryptMsg(echostr, srv.EncodingAESKey)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Fprintf(srv.Writer, string(xml_content))

	}
	if srv.Request.Method == "POST" {
		fmt.Println("it is post")
		/*
			response, err := srv.handleRequest()
			if err != nil {
				return err
			}

			debug
			fmt.Println("request msg = ", string(srv.requestRawXMLMsg))

			return srv.buildResponse(response)
		*/
	}
}

//Validate 校验消息体请求是否合法
func (srv *Server) Validate(str string) bool {
	timestamp := srv.GetURLParam("timestamp")
	nonce := srv.GetURLParam("nonce")
	signature := srv.GetURLParam("msg_signature")
	return msgcrypt.ValidateMsg(srv.Token, timestamp, nonce, str, signature)
}

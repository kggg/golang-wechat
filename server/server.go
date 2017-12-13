package server

import (
	"fmt"
	//"io/ioutil"
	"encoding/xml"
	"log"
	"strconv"
	"strings"
	"time"
	"wechat/context"
	"wechat/message"
	"wechat/monitor"
	"wechat/utils/msgcrypt"
)

//Server struct
type Server struct {
	*context.Context
	random      []byte
	msg_encrypt string
	nonce       string
	timestamp   int64
}

//NewServer init
func NewServer(context *context.Context) *Server {
	srv := new(Server)
	srv.Context = context
	return srv
}

//Serve 处理微信的请求消息
func (srv *Server) Serve() error {
	//显示请求URL
	log.Println(srv.Request.Method, srv.Request.URL.RequestURI())
	//验证请求体是否合法
	if !srv.Validate() {
		return fmt.Errorf("消息不合法，验证签名失败")
	}
	//微信URL回调模式验证
	if srv.Request.Method == "GET" {
		xml_content, err := msgcrypt.DecryptMsg(srv.msg_encrypt, srv.EncodingAESKey)
		if err != nil {
			log.Fatalln(err)
			return fmt.Errorf("解密xml失败")
		}
		fmt.Fprintf(srv.Writer, string(xml_content))
		return nil

	}
	//接收消息处理
	if srv.Request.Method == "POST" {
		//Get Raw post xml info and decrypt
		req_content, err := msgcrypt.DecryptMsg(srv.msg_encrypt, srv.EncodingAESKey)
		if err != nil {
			log.Fatalln(err)
			return fmt.Errorf("解密xml失败")
		}
		//parse plaintext xml info
		text := &message.Text{}
		err = xml.Unmarshal(req_content, &text)
		if err != nil {
			fmt.Fprintf(srv.Writer, "error: %v", err)
			return err
		}
		msgtype := text.MsgType
		switch msgtype {
		case "text":
		case "image":
		case "voice":
		case "video":
		case "music":
		case "news":
		case "transfer":
		default:
			err = message.ErrUnsupportReply
			fmt.Fprintf(srv.Writer, "不支持 MsgTyp=%s 的消息格式", msgtype)
			return err
		}
		//被动响应消息, 这里用来测试内部主机的状态
		if strings.HasPrefix(text.Content, "status") {
			var msg string
			content := strings.Split(text.Content, "#")
			if len(content) > 2 {
				if content[1] == "host" {
					msg = "the host status is active"
					ok := monitor.Ping(content[2])
					if !ok {
						msg = "the host status is down"
					}
				} else if content[1] == "service" {

				} else {
					msg = "你输入格式不对\n格式为: status#[host|service]#[ip|servicename]\n"
				}
			} else {
				msg = "你输入格式不对\n格式为: status#[host|service]#[ip|servicename]\n"
			}

			responsexml, err := srv.MakeResponseMsg(msg, text.FromUserName, text.MsgType)
			if err != nil {
				return err
			}
			fmt.Fprintf(srv.Writer, string(responsexml))
		}
		//Response content
		return nil
	}
	return nil
}

func (srv *Server) MakeResponseMsg(msg, touser, msgtype string) ([]byte, error) {
	text := message.NewText(msg)
	text.ToUserName = touser
	text.CreateTime = time.Now().Unix()
	text.MsgType = msgtype
	text.FromUserName = srv.AppID
	textmsg, err := xml.MarshalIndent(text, "", "    ")
	if err != nil {
		return nil, err
	}
	cipher, err := msgcrypt.EncryptMsg(string(textmsg), srv.EncodingAESKey, srv.AppID)
	if err != nil {
		return nil, err
	}
	timesp := strconv.FormatInt(text.CreateTime, 10)
	signature := msgcrypt.MakeSHA1Slice(srv.Token, timesp, srv.nonce, cipher)
	body, err := message.MakeResponseXML(cipher, signature, srv.nonce, text.CreateTime)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//Validate 校验消息体请求是否合法
func (srv *Server) Validate() bool {
	timestamp := srv.GetURLParam("timestamp")
	times, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	srv.timestamp = times
	srv.nonce = srv.GetURLParam("nonce")
	echostr, exists := srv.QueryURLParam("echostr")
	if !exists {
		var encryptxml message.EncryptedXMLMsg
		encryptxml, err := srv.GetRequestBody()
		if err != nil {
			return false
		}
		echostr = encryptxml.EncryptedMsg
	}
	srv.msg_encrypt = echostr
	signature := srv.GetURLParam("msg_signature")
	return msgcrypt.ValidateMsg(srv.Token, timestamp, srv.nonce, echostr, signature)
}

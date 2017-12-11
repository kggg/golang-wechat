package message

import (
	"encoding/xml"
)

// MsgType 基本消息类型
var MsgType = []string{"text", "image", "voice", "video", "shortvideo", "location", "link", "music", "news", "transfer_customer_service", "event"}

type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

//EncryptedXMLMsg 安全模式下的消息体
type EncryptedXMLMsg struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	ToUserName   string   `xml:"ToUserName" json:"ToUserName"`
	EncryptedMsg string   `xml:"Encrypt"    json:"Encrypt"`
	AgentID      int      `xml:"AgentID"    json:"AgentID"`
}

//ResponseEncryptedXMLMsg 需要返回的消息体
type ResponseEncryptedXMLMsg struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	Encrypt      CDATA    `json:"Encrypt"`
	MsgSignature CDATA    `json:"MsgSignature"`
	TimeStamp    int64    `xml:"TimeStamp"  json:"TimeStamp"`
	Nonce        CDATA    `json:"Nonce"`
}

// CommonToken 消息中通用的结构
type CommonToken struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
}

//SetToUserName set ToUserName
func (msg *CommonToken) SetToUserName(toUserName string) {
	msg.ToUserName = toUserName
}

//SetFromUserName set FromUserName
func (msg *CommonToken) SetFromUserName(fromUserName string) {
	msg.FromUserName = fromUserName
}

//SetCreateTime set createTime
func (msg *CommonToken) SetCreateTime(createTime int64) {
	msg.CreateTime = createTime
}

//SetMsgType set MsgType
func (msg *CommonToken) SetMsgType(msgType string) {
	msg.MsgType = msgType
}

type FullMessage struct {
	CommonToken

	//基本消息
	MsgID        int64   `xml:"MsgId"`
	Content      string  `xml:"Content"`
	PicURL       string  `xml:"PicUrl"`
	MediaID      string  `xml:"MediaId"`
	Format       string  `xml:"Format"`
	ThumbMediaID string  `xml:"ThumbMediaId"`
	LocationX    float64 `xml:"Location_X"`
	LocationY    float64 `xml:"Location_Y"`
	Scale        float64 `xml:"Scale"`
	Label        string  `xml:"Label"`
	Title        string  `xml:"Title"`
	Description  string  `xml:"Description"`
	URL          string  `xml:"Url"`

	//事件相关
	Event     string `xml:"Event"`
	EventKey  string `xml:"EventKey"`
	Ticket    string `xml:"Ticket"`
	Latitude  string `xml:"Latitude"`
	Longitude string `xml:"Longitude"`
	Precision string `xml:"Precision"`
	MenuID    string `xml:"MenuId"`
}

func MakeResponseXML(cipher, signature, nonce string, timestamp int64) ([]byte, error) {
	/*
		msgXML := &{
			Encrypt:      CDATA(cipher),
			MsgSignature: CDATA(signature),
			TimeStamp:    timestamp,
			Nonce:        CDATA(nonce),
		}
	*/
	var resp ResponseEncryptedXMLMsg
	resp.Encrypt = CDATA(cipher)
	resp.MsgSignature = CDATA(signature)
	resp.TimeStamp = timestamp
	resp.Nonce = CDATA(nonce)
	body, err := xml.MarshalIndent(resp, " ", "    ")
	if err != nil {
		return nil, err
	}
	return body, nil
}

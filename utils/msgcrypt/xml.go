package msgcrypt

import (
	"encoding/xml"
)

type ResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATA
	MsgSignature CDATA
	TimeStamp    int64
	Nonce        CDATA
}

type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

type RequestBody struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	AgentID    string
	Encrypt    string
}

type CommonMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
}

type MsgBody struct {
	CommonMsg
	Content string `xml:"Content"`
}

func ParseEncryptXml(req RequestBody, key string) {
	body := req.Body
	request, err := ParseXml(body)
	xml_content, err := DecryptMsg(request.Encrypt, key)
	if err != nil {
		log.Fatalln(err)
	}
	var rec MsgBody
	err = xml.Unmarshal(xml_content, &rec)
	if err != nil {
		beego.Debug("msgbody invalid:", err)
		log.Fatalln(err)
	}
}

func ParseXml(xml string, xmlStruct interface{}) (interface{}, error) {
	var request xmlStruct
	err := xml.Unmarshal(body, &request)
	return request, err
}

package context

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"wechat/message"
)

// Context struct
type Context struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string

	//Cache cache.Cache

	Writer  http.ResponseWriter
	Request *http.Request

	//accessTokenLock 读写锁 同一个AppID一个
	//accessTokenLock *sync.RWMutex

	//jsAPITicket 读写锁 同一个AppID一个
	//jsAPITicketLock *sync.RWMutex
}

// Query returns the keyed url query value if it exists
func (ctx *Context) GetURLParam(key string) string {
	value, _ := ctx.QueryURLParam(key)
	return value
}

// GetQuery is like Query(), it returns the keyed url query value
func (ctx *Context) QueryURLParam(key string) (string, bool) {
	values, ok := ctx.Request.URL.Query()[key]
	if ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func (ctx *Context) GetRequestBody() (message.EncryptedXMLMsg, error) {
	var encryptxml message.EncryptedXMLMsg
	if err := xml.NewDecoder(ctx.Request.Body).Decode(&encryptxml); err != nil {
		return encryptxml, fmt.Errorf("从body中解析xml失败,err=%v", err)
	}
	return encryptxml, nil
}

package msgcrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

func MakeSHA1Slice(token, timestamp, nonce, msg_encrypt string) string {
	sl := []string{token, timestamp, nonce, msg_encrypt}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func ValidateMsg(token, timestamp, nonce, msg_encrypt, msgSignature string) bool {
	checkmsg := MakeSHA1Slice(token, timestamp, nonce, msg_encrypt)
	if checkmsg != msgSignature {
		return false
	}
	return true
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length < 1 {
		return nil, errors.New("the cipther length less than 1")
	}
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)], nil
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func EncryptMsg(plaintMsg, key, corpid string) (string, error) {
	aeskey, err := base64.StdEncoding.DecodeString(key + "=")
	if err != nil {
		return "", err
	}
	plaintext, err := JoinMsg(plaintMsg, corpid)
	cipherMsg, err := AesEncrypt(plaintext, aeskey)
	if err != nil {
		return "", err
	}
	base64msg := base64.StdEncoding.EncodeToString(cipherMsg)
	if err != nil {
		return "", err
	}
	return base64msg, nil
}

func JoinMsg(msg, corpid string) ([]byte, error) {
	randomBytes := []byte("abcdefghabcdefgh")
	Msg := []byte(msg)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(msg)))
	if err != nil {
		return []byte(""), err
	}
	msgLength := buf.Bytes()
	plaintext := bytes.Join([][]byte{randomBytes, msgLength, Msg, []byte(corpid)}, nil)
	return plaintext, nil

}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

func DecryptMsg(EncryptMsg string, key string) ([]byte, error) {
	aeskey, err := base64.StdEncoding.DecodeString(key + "=")
	if err != nil {
		return nil, err
	}
	aes_msg, err := base64.StdEncoding.DecodeString(EncryptMsg)
	if err != nil {
		return nil, err
	}
	plaintext, err := AesDecrypt(aes_msg, aeskey)
	if err != nil {
		return nil, err
	}
	bytesBuffer := bytes.NewBuffer(plaintext[16:20])
	var msg_int int32
	binary.Read(bytesBuffer, binary.BigEndian, &msg_int)
	msgint := int(msg_int)
	xml_content := plaintext[20 : 20+msgint]
	return xml_content, nil
}

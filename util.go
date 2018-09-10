package wxpay

import (
	"bytes"
	"crypto/tls"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/pkcs12"
	"log"
	"strings"
)

func XmlToMap(xmlStr string) Params {

	params := make(Params)
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	var (
		key   string
		value string
	)

	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 开始标签
			key = token.Name.Local
		case xml.CharData: // 标签内容
			content := string([]byte(token))
			value = content
		}
		if key != "xml" {
			if value != "\n" {
				params.SetString(key, value)
			}
		}
	}

	return params
}

func MapToXml(params Params) string {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`><![CDATA[`)
		buf.WriteString(v)
		buf.WriteString(`]]></`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)

	return buf.String()
}

// 生成随机字符串
func nonceStr() string {
	uid, err := uuid.NewV4()
	if err != nil {
		return errors.New("生成随机字符串").Error()
	}
	return strings.Replace(uid.String(), "-", "", -1)
}

// 将Pkcs12转成Pem
func pkcs12ToPem(p12 []byte, password string) tls.Certificate {

	blocks, err := pkcs12.ToPEM(p12, password)

	// 从恐慌恢复
	defer func() {
		if x := recover(); x != nil {
			log.Print(x)
		}
	}()

	if err != nil {
		panic(err)
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		panic(err)
	}
	return cert
}

package wxpay

import (
	"bytes"
	"encoding/xml"
	"errors"
	"github.com/satori/go.uuid"
	"io"
	"strings"
)

func XmlToMap(r io.Reader) Params {

	params := make(Params)

	decoder := xml.NewDecoder(r)

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

func MapToXml(reqData Params) string {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range reqData {
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
func NonceStr() string {
	uid, err := uuid.NewV4()
	if err != nil {
		return errors.New("生成随机字符串").Error()
	}
	return strings.Replace(uid.String(), "-", "", -1)
}

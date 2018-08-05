package wxpay

import (
	"bytes"
	"encoding/xml"
	"errors"
	"github.com/satori/go.uuid"
	"strings"
)

func XmlToMap(strXml string) Params {

	params := make(Params)

	inputReader := strings.NewReader(strXml)
	decoder := xml.NewDecoder(inputReader)

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
			params.SetString(key, value)
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

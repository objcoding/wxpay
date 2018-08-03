package wxpay

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
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
		params.SetString(key, value)
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

func Sign(params Params) string {
	// 创建切片
	var keys = make([]string, 0, len(params))
	// 遍历签名参数
	for k, _ := range params {
		if k != "sign" { // 排除sign字段
			keys = append(keys, k)
		}
	}

	//创建字符缓冲
	var buf bytes.Buffer
	for _, k := range keys {
		if len(params.GetString(k)) > 0 {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(params.GetString(k))
			buf.WriteString(`&`)
		}
	}
	// 加入apiKey作加密密钥
	buf.WriteString(`key=`)
	buf.WriteString(config.Key)

	var (
		dataMd5    [16]byte
		dataSha256 [32]byte
		str        string
	)

	switch params.GetString("sign_type") {
	case MD5:
		dataMd5 = md5.Sum(buf.Bytes())
		str = hex.EncodeToString(dataMd5[:]) //需转换成切片
	case HMACSHA256:
		dataSha256 = sha256.Sum256(buf.Bytes())
		str = hex.EncodeToString(dataSha256[:])
	}

	return strings.ToUpper(str)
}

func GenerateSignedXml(params Params) string {
	sign := Sign(params)
	params.SetString(FIELD_SIGN, sign)
	return MapToXml(params)
}

func IsSignatureValid(xmlStr string) bool {
	params := XmlToMap(xmlStr)
	if !params.ContainsKey(FIELD_SIGN) {
		return false
	}
	sign := params.GetString(FIELD_SIGN)
	newSign := Sign(params)
	return sign == newSign
}

// 生成随机字符串
func NonceStr() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", errors.New("生成随机字符串")
	}
	return strings.Replace(uid.String(), "-", "", -1), nil
}

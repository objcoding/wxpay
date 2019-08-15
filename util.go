package wxpay

import (
	"crypto/tls"
	"encoding/pem"
	"encoding/xml"
	"golang.org/x/crypto/pkcs12"
	"log"
	"strconv"
	"time"
	"io"
	"bytes"
)

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string 	`xml:",chardata"`
}

func (m Params) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}
	start.Name.Local = "xml" // 更改xml开始标签
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func (m *Params) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Params{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

func XmlToMap(xmlStr string) Params {
	params := make(Params)
	err := xml.Unmarshal([]byte(xmlStr), &params)
	if err != nil {
		panic(err)
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

// 用时间戳生成随机字符串
func nonceStr() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
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

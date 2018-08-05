package wxpay

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

const bodyType = "application/xml; charset=utf-8"

type Client struct {
	account    *Account
	signType   string
	useSandbox bool
	httpClient *http.Client
}

func NewClient(account *Account, signType string, useSandbox bool) *Client {
	return &Client{
		account:    account,
		signType:   signType,
		useSandbox: useSandbox,
	}
}

//向 Map 中添加 appid、mch_id、nonce_str、sign_type、sign
// 该函数适用于商户适用于统一下单等接口，不适用于红包、代金券接口
func (c *Client) FillRequestData(params Params) Params {
	params["appid"] = c.account.AppID
	params["mch_id"] = c.account.MchID
	params["nonce_str"] = NonceStr()
	params["sign_type"] = MD5
	params["sign"] = c.Sign(params)

	return params
}

func (c *Client) Post(url string, params Params) {

}

func (c *Client) GenerateSignedXml(params Params) string {
	sign := c.Sign(params)
	params.SetString(FIELD_SIGN, sign)
	return MapToXml(params)
}

func (c *Client) ValidSign(xmlStr string) bool {
	params := XmlToMap(xmlStr)
	if !params.ContainsKey(FIELD_SIGN) {
		return false
	}
	return params.GetString(FIELD_SIGN) == c.Sign(params)
}

// 签名
func (c *Client) Sign(params Params) string {
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
	buf.WriteString(c.account.Key)

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

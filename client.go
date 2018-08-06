package wxpay

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

const bodyType = "application/xml; charset=utf-8"

type Client struct {
	account              *Account // 支付账号
	signType             string   // 签名类型
	useSandbox           bool     // 是否为沙箱
	httpConnectTimeoutMs int      // 连接超时时间
	httpReadTimeoutMs    int      // 读取超时时间
}

// 创建微信支付客户端
func NewClient(account *Account, useSandbox bool) *Client {
	return &Client{
		account:              account,
		signType:             MD5,
		useSandbox:           useSandbox,
		httpConnectTimeoutMs: 2000,
		httpReadTimeoutMs:    1000,
	}
}

func (c *Client) SetHttpConnectTimeoutMs(ms int) {
	c.httpConnectTimeoutMs = ms
}

func (c *Client) SetHttpReadTimeoutMs(ms int) {
	c.httpReadTimeoutMs = ms
}

func (c *Client) SetSignType(signType string) {
	c.signType = signType
}

func (c *Client) setAccount(account *Account) {
	c.account = account
}

// 向 params 中添加 appid、mch_id、nonce_str、sign_type、sign
func (c *Client) fillRequestData(params Params) Params {
	params["appid"] = c.account.AppID
	params["mch_id"] = c.account.MchID
	params["nonce_str"] = NonceStr()
	params["sign_type"] = c.signType
	params["sign"] = c.Sign(params)
	return params
}

// https no cert post
func (c *Client) postWithoutCert(url string, params Params) (string, error) {
	h := &http.Client{}
	p := c.fillRequestData(params)
	response, err := h.Post(url, bodyType, strings.NewReader(MapToXml(p)))
	if err != nil {
		return "", err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// https need cert post
func (c *Client) postWithCert(url string, params Params) (string, error) {
	if c.account.CertData == nil {
		return "", errors.New("证书数据为空")
	}
	_, certificate, err := pkcs12.Decode(c.account.CertData, c.account.MchID)
	if err != nil {
		return "", err
	}
	pool := x509.NewCertPool()
	pool.AddCert(certificate)

	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: true,
	}
	h := &http.Client{Transport: transport}
	p := c.fillRequestData(params)
	response, err := h.Post(url, bodyType, strings.NewReader(MapToXml(p)))
	if err != nil {
		return "", err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// 生成带有签名的xml字符串
func (c *Client) GenerateSignedXml(params Params) string {
	sign := c.Sign(params)
	params.SetString(FIELD_SIGN, sign)
	return MapToXml(params)
}

// 验证签名
func (c *Client) ValidSign(params Params) bool {
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
	for k := range params {
		if k != "sign" { // 排除sign字段
			keys = append(keys, k)
		}
	}

	// 由于切片的元素顺序是不固定，所以这里强制给切片元素加个顺序
	sort.Strings(keys)

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
	buf.WriteString(c.account.ApiKey)

	var (
		dataMd5    [16]byte
		dataSha256 [32]byte
		str        string
	)

	switch c.signType {
	case MD5:
		dataMd5 = md5.Sum(buf.Bytes())
		str = hex.EncodeToString(dataMd5[:]) //需转换成切片
	case HMACSHA256:
		dataSha256 = sha256.Sum256(buf.Bytes())
		str = hex.EncodeToString(dataSha256[:])
	}

	return strings.ToUpper(str)
}

// 处理 HTTPS API返回数据，转换成Map对象。return_code为SUCCESS时，验证签名。
func (c *Client) processResponseXml(xmlStr string) (Params, error) {
	var returnCode string
	params := XmlToMap(strings.NewReader(xmlStr))
	if params.ContainsKey("return_code") {
		returnCode = params.GetString("return_code")
	} else {
		return nil, errors.New("no return_code in XML")
	}
	if returnCode == FAIL {
		return params, nil
	} else if returnCode == SUCCESS {
		if c.ValidSign(params) {
			return params, nil
		} else {
			return nil, errors.New("invalid sign value in XML")
		}
	} else {
		return nil, errors.New("return_code value is invalid in XML")
	}
}

// 统一下单
func (c *Client) UnifiedOrder(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_UNIFIEDORDER_URL
	} else {
		url = UNIFIEDORDER_URL
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

// 刷卡支付
func (c *Client) MicroPay(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_MICROPAY_URL
	} else {
		url = MICROPAY_URL
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

// 退款
func (c *Client) Refund(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_REFUND_URL
	} else {
		url = REFUND_URL
	}
	xmlStr, err := c.postWithCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

// 订单查询
func (c *Client) OrderQuery(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_ORDERQUERY_URL
	} else {
		url = ORDERQUERY_URL
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

// 退款查询
func (c *Client) RefundQuery(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_REFUNDQUERY_URL
	} else {
		url = REFUNDQUERY_URL
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

// 撤销订单
func (c *Client) Reverse(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_REVERSE_URL
	} else {
		url = REVERSE_URL
	}
	xmlStr, err := c.postWithCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

// 关闭订单
func (c *Client) CloseOrder(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_CLOSEORDER_URL
	} else {
		url = CLOSEORDER_URL
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

// 对账单下载
func (c *Client) DownloadBill(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_DOWNLOADBILL_URL
	} else {
		url = DOWNLOADBILL_URL
	}
	xmlStr, err := c.postWithoutCert(url, params)

	var p Params

	// 如果出现错误，返回XML数据
	if strings.Index(xmlStr, "<") == 0 {
		p = XmlToMap(strings.NewReader(xmlStr))
		return p, err
	} else { // 正常返回csv数据
		p.SetString("return_code", SUCCESS)
		p.SetString("return_msg", "ok")
		p.SetString("data", xmlStr)
		return p, err
	}
}

// 交易保障
func (c *Client) Report(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_REPORT_URL
	} else {
		url = REPORT_URL
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

// 转换短链接
func (c *Client) ShortUrl(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_SHORTURL_URL
	} else {
		url = SHORTURL_URL
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

// 授权码查询OPENID接口
func (c *Client) AuthCodeToOpenid(params Params) (Params, error) {
	var url string
	if c.useSandbox {
		url = SANDBOX_AUTHCODETOOPENID_URL
	} else {
		url = AUTHCODETOOPENID_URL
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

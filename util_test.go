package wxpay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXmlToMap(t *testing.T) {
	xmlStr := "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg><appid><![CDATA[wx2421b1c4370ec43b]]></appid><mch_id><![CDATA[10000100]]></mch_id><nonce_str><![CDATA[IITRi8Iabbblz1Jc]]></nonce_str><sign><![CDATA[7921E432F65EB8ED0CE9755F0E86D72F]]></sign><result_code><![CDATA[SUCCESS]]></result_code><prepay_id><![CDATA[wx201411101639507cbf6ffd8b0779950874]]></prepay_id><trade_type><![CDATA[APP]]></trade_type></xml>"
	params := XmlToMap(xmlStr)
	if params == nil {
		t.Error(params)
	}
	t.Log(params)
}

func TestXmlToMapWithLineBreaksAndSpaces(t *testing.T) {
	xmlStr := `<xml>
  <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
  <attach><![CDATA[支付测试]]></attach>
  <bank_type><![CDATA[CFT]]></bank_type>
  <fee_type><![CDATA[CNY]]></fee_type>
  <is_subscribe><![CDATA[Y]]></is_subscribe>
  <mch_id><![CDATA[10000100]]></mch_id>
  <nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
  <openid><![CDATA[oUpF8uMEb4qRXf22hE3X68TekukE]]></openid>
  <out_trade_no><![CDATA[1409811653]]></out_trade_no>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <sign><![CDATA[B552ED6B279343CB493C5DD0D78AB241]]></sign>
  <time_end><![CDATA[20140903131540]]></time_end>
  <total_fee>1</total_fee>
  <coupon_fee><![CDATA[10]]></coupon_fee>
  <coupon_count><![CDATA[1]]></coupon_count>
  <coupon_type><![CDATA[CASH]]></coupon_type>
  <coupon_id><![CDATA[10000]]></coupon_id>
  <trade_type><![CDATA[JSAPI]]></trade_type>
  <transaction_id><![CDATA[1004400740201409030005092168]]></transaction_id>
</xml>
`
	params := XmlToMap(xmlStr)
	if params == nil {
		t.Error(params)
	}
	t.Log(params)

	assert.Equal(t, "wx2421b1c4370ec43b", params["appid"])
	assert.Equal(t, "支付测试", params["attach"])
	assert.Equal(t, "1004400740201409030005092168", params["transaction_id"])
}

func TestMapToXml(t *testing.T) {
	params := map[string]string{"return_msg": "OK", "appid": "wx2421b1c4370ec43b", "mch_id": "10000100"}
	xmlStr := MapToXml(params)
	t.Log(xmlStr)
}

func TestNonceStr(t *testing.T) {
	t.Log(nonceStr())
}

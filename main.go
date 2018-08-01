package main

import (
	"fmt"
	"go-pay/wxpay"
)

func main() {

	//strXml := "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg><appid><![CDATA[wx2421b1c4370ec43b]]></appid><mch_id><![CDATA[10000100]]></mch_id><nonce_str><![CDATA[NfsMFbUFpdbEhPXP]]></nonce_str></xml>"
	//strXml2 := "<return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg>"

	params := &wxpay.Params{
		"return_code": "SUCCESS",
		"return_msg":  "OK",
	}

	fmt.Print(wxpay.MapToXml(params))
}

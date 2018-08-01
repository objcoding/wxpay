package wxpay

type Client struct {
	config     Config
	signType   string
	useSandbox bool
}

func (w *Client) Default(config Config) *Client {
	return &Client{
		config:     config,
		signType:   MD5,
		useSandbox: false,
	}
}

func (w *Client) NewClient(config Config, signType string, useSandbox bool) *Client {
	return &Client{
		config:     config,
		signType:   signType,
		useSandbox: useSandbox,
	}
}

//向 Map 中添加 appid、mch_id、nonce_str、sign_type、sign
// 该函数适用于商户适用于统一下单等接口，不适用于红包、代金券接口
func FillRequestData(reqData map[string]string) map[string]string {
	reqData["appid"] = config.AppID
	reqData["mch_id"] = config.MchID
	reqData["nonce_str"] = ""
	reqData["sign_type"] = "MD5"
	reqData["sign"] = ""

	return reqData
}

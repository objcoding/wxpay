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

func (w *Client) New(config Config, signType string, useSandbox bool) *Client {
	return &Client{
		config:     config,
		signType:   signType,
		useSandbox: useSandbox,
	}
}

//向 Map 中添加 appid、mch_id、nonce_str、sign_type、sign
// 该函数适用于商户适用于统一下单等接口，不适用于红包、代金券接口
func FillRequestData(params Params) map[string]string {
	params["appid"] = config.AppID
	params["mch_id"] = config.MchID
	params["nonce_str"] = ""
	params["sign_type"] = "MD5"
	params["sign"] = Sign(params)

	return params
}

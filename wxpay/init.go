package wxpay

var config *Config

// 初始化
func init() {
	config = Instance("appid", "mchId", "key", 200, 500)
}

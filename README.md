# wxpay for golang

## 微信支付开发文档API

对[微信支付开发者文档](https://pay.weixin.qq.com/wiki/doc/api/index.html)中给出的API进行了封装。

wxpay提供了对应的方法：

| 方法名              | 说明          |
| ---------------- | ----------- |
| microPay         | 刷卡支付        |
| unifiedOrder     | 统一下单        |
| orderQuery       | 查询订单        |
| reverse          | 撤销订单        |
| closeOrder       | 关闭订单        |
| refund           | 申请退款        |
| refundQuery      | 查询退款        |
| downloadBill     | 下载对账单       |
| report           | 交易保障        |
| shortUrl         | 转换短链接       |
| authCodeToOpenid | 授权码查询openid |


## 安装

```
go get github.com/objcoding/wxpay

```


## 示例

```go

// 新建微信支付客户端
client := wxpay.NewClient(wxpay.NewAccount{
	AppID: "appid",
	MchID: "mchid",
	ApiKey: "apiKey",
}, false) // sandbox环境请传true

// 统一下单
params := make(wxpay.Params)
	params.SetString("body", "test").
		SetString("out_trade_no", "436577857").
		SetInt64("total_fee", 1).
		SetString("spbill_create_ip", "127.0.0.1").
		SetString("notify_url", "http://objcoding.com").
		SetString("trade_type", "APP")
p, _ := client.UnifiedOrder(params)

// 订单查询
params := make(wxpay.Params)
params.SetString("out_trade_no", "3568785")
p, _ := client.OrderQuery(params)

// 退款
params := make(wxpay.Params)
params.SetString("body", "test").
    SetString("out_trade_no", "3568785").
    SetInt64("total_fee", 1).
    SetInt64("refund_fee", 1)
p, _ := client.Refund(params)

// 退款查询
params := make(wxpay.Params)
params.SetString("out_refund_no", "3568785")
p, _ := client.RefundQuery(params)

```


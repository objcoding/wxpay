# wxpay for golang

注：本人照着微信支付官方 SDK Java 版，用 Go 语言实现的一个微信支付 SDK，与 Java 版功能一致，具备轻量易上手，且功能齐全的特点，理论上性能也会比 Java 版的更快。

## 微信支付开发文档API

对[微信支付开发者文档](https://pay.weixin.qq.com/wiki/doc/api/index.html)中给出的API进行了封装。

wxpay提供了对应的方法：

| 方法名              | 说明          |
| ---------------- | ----------- |
| MicroPay         | 刷卡支付        |
| UnifiedOrder     | 统一下单        |
| OrderQuery       | 查询订单        |
| Reverse          | 撤销订单        |
| CloseOrder       | 关闭订单        |
| Refund           | 申请退款        |
| RefundQuery      | 查询退款        |
| DownloadBill     | 下载对账单       |
| Report           | 交易保障        |
| ShortUrl         | 转换短链接       |
| AuthCodeToOpenid | 授权码查询openid |

* 参数为`Map<String, String>`对象，返回类型也是`Map<String, String>`。
* 方法内部会将参数会转换成含有`appid`、`mch_id`、`nonce_str`、`sign\_type`和`sign`的XML；
* 默认使用MD5进行签名；
* 通过HTTPS请求得到返回数据后会对其做必要的处理（例如验证签名，签名错误则抛出异常）。
* 对于downloadBill，无论是否成功都返回Map，且都含有`return_code`和`return_msg`。若成功，其中`return_code`为`SUCCESS`，另外`data`对应对账单数据。


## 安装

```
go get github.com/objcoding/wxpay

```


## 示例

```cgo

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

```cgo

// 创建支付账户
account := wxpay.NewAccount("appid", "mchid", "apiKey")

// 设置证书
account.SetCertData("证书地址")

// 新建微信支付客户端
client := wxpay.NewClient(account, false) // sandbox环境请传true

// 设置http请求超时时间
client.SetHttpConnectTimeoutMs(2000)

// 设置http读取信息流超时时间
client.SetHttpReadTimeoutMs(1000)

// 更改签名类型
client.SetSignType(HMACSHA256)

// 设置支付账户
client.setAccount(account)

```

```cgo

// 签名
signStr := client.Sign(params)

// 校验签名
b := client.ValidSign(params)

```

```cgo

// 支付或退款返回成功信息
return wxpay.Notifies{}.OK()

// 支付或退款返回失败信息
wxpay.Notifies{}.NotOK("支付失败或退款失败了")

```


## License
BSD

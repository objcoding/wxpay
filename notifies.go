package wxpay

type Notifies struct{}

// 通知成功
func (n *Notifies) OK() string {
	var params = make(Params)
	params.SetString("return_code", Success)
	params.SetString("return_msg", "ok")
	return MapToXml(params)
}

// 通知不成功
func (n *Notifies) NotOK(errMsg string) string {
	var params = make(Params)
	params.SetString("return_code", Fail)
	params.SetString("return_msg", errMsg)
	return MapToXml(params)
}

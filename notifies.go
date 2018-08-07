package wxpay

type Notifies struct{}

// 通知成功
func (n *Notifies) OK() string {
	var params = make(Params)
	params.SetString("return_code", SUCCESS)
	params.SetString("return_msg", "ok")
	return mapToXml(params)
}

// 通知不成功
func (n *Notifies) NotOK(errMsg string) string {
	var params = make(Params)
	params.SetString("return_code", FAIL)
	params.SetString("return_msg", errMsg)
	return mapToXml(params)
}

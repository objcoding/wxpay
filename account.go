package wxpay

import (
	"io/ioutil"
	"log"
)

type Account struct {
	appID     string
	subAppId  string
	mchID     string
	subMchId  string
	apiKey    string
	certData  []byte
	isSandbox bool
}

// 创建微信支付账号
func NewAccount(appID string, subAppId string, mchID string, subMchId string, apiKey string, isSandbox bool) *Account {
	return &Account{
		appID:     appID,
		subAppId:  subAppId,
		mchID:     mchID,
		subMchId:  subMchId,
		apiKey:    apiKey,
		isSandbox: isSandbox,
	}
}

// 设置证书
func (a *Account) SetCertData(certPath string) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Println("读取证书失败")
		return
	}
	a.certData = certData
}

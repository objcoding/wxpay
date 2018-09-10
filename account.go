package wxpay

import (
	"io/ioutil"
	"log"
)

type Account struct {
	appID     string
	mchID     string
	apiKey    string
	certData  []byte
	isSandbox bool
}

// 创建微信支付账号
func NewAccount(appID string, mchID string, apiKey string, isSanbox bool) *Account {
	return &Account{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		isSandbox: isSanbox,
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

package wxpay

import (
	"io/ioutil"
	"log"
)

type Account struct {
	AppID     string
	MchID     string
	ApiKey    string
	CertData  []byte
	isSandbox bool
}

// 创建微信支付账号
func NewAccount(AppID string, MchID string, ApiKey string, isSanbox bool) *Account {
	return &Account{
		AppID:     AppID,
		MchID:     MchID,
		ApiKey:    ApiKey,
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
	a.CertData = certData
}

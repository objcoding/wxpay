package wxpay

import (
	"io/ioutil"
	"log"
)

type Account struct {
	AppID    string
	MchID    string
	ApiKey   string
	CertData []byte
}

// 创建微信支付账号
func NewAccount(AppID string, MchID string, ApiKey string) *Account {
	return &Account{
		AppID:  AppID,
		MchID:  MchID,
		ApiKey: ApiKey,
	}
}

// 设置证书
func (a *Account) setCertData(certPath string) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Println("读取证书失败")
		return
	}
	a.CertData = certData
}

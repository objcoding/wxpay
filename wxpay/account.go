package wxpay

import "io"

type Account struct {
	AppID                string
	MchID                string
	Key                  string
	certStream           io.Reader
	HttpConnectTimeoutMs int
	HttpReadTimeoutMs    int
}

func NewConfig(AppID string, MchID string, Key string, HttpConnectTimeoutMs int, HttpReadTimeoutMs int) *Account {
	return &Account{
		AppID:                AppID,
		MchID:                MchID,
		Key:                  Key,
		HttpConnectTimeoutMs: HttpConnectTimeoutMs,
		HttpReadTimeoutMs:    HttpReadTimeoutMs,
	}
}

func (a *Account) setCertStream(reader io.Reader) {
	a.certStream = reader
}

package wxpay

import "io"

type Config struct {
	AppID                string
	MchID                string
	Key                  string
	certStream           io.Reader
	HttpConnectTimeoutMs int
	HttpReadTimeoutMs    int
}

func Instance(AppID string, MchID string, Key string, HttpConnectTimeoutMs int, HttpReadTimeoutMs int) *Config {
	return &Config{
		AppID:                AppID,
		MchID:                MchID,
		Key:                  Key,
		HttpConnectTimeoutMs: HttpConnectTimeoutMs,
		HttpReadTimeoutMs:    HttpReadTimeoutMs,
	}
}

func (c *Config) setCertStream(reader io.Reader) {
	c.certStream = reader
}

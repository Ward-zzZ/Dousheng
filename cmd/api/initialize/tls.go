package initialize

import (
	"crypto/tls"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// InitTLS
// 自签名证书生成：openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
func InitTLS() *tls.Config {
	cfg := &tls.Config{
		MinVersion:         tls.VersionTLS10,
		InsecureSkipVerify: true,
	}
	cert, err := tls.LoadX509KeyPair("./cert/server.crt",
		"./cert/server.key")
	if err != nil {
		hlog.Fatal("tls failed", err)
	}
	cfg.Certificates = append(cfg.Certificates, cert)
	return cfg
}

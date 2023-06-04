package initialize

import (
	"net"
	"strconv"

	"tiktok-demo/cmd/relation/config"
	"tiktok-demo/shared/consts"
	"tiktok-demo/shared/tools"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// InitConfig to init consul config server
func InitConfig() {
	v := viper.New()
	v.SetConfigFile(consts.ConfigPath)
	if err := v.ReadInConfig(); err != nil {
		hlog.Fatal("read config file failed: %s", err.Error())
	}
	if err := v.Unmarshal(&config.GlobalConsulConfig); err != nil {
		hlog.Fatal("unmarshal config failed: %s", err.Error())
	}
	hlog.Info("Config Info: %v", config.GlobalConsulConfig)

	cfg := api.DefaultConfig()
	cfg.Address = net.JoinHostPort(
		config.GlobalConsulConfig.Host,
		strconv.Itoa(config.GlobalConsulConfig.Port))
	consulClient, err := api.NewClient(cfg)
	if err != nil {
		hlog.Fatal("create consul client failed: %s", err.Error())
	}
	content, _, err := consulClient.KV().Get(config.GlobalConsulConfig.Key, nil)
	if err != nil {
		hlog.Fatal("consul KV get failed: %s", err.Error())
	}

	err = sonic.Unmarshal(content.Value, &config.GlobalServerConfig)
	if err != nil {
		hlog.Fatal("sonic unmarshal config failed: %s", err.Error())
	}

	if config.GlobalServerConfig.Host == "" {
		config.GlobalServerConfig.Host, err = tools.GetLocalIPv4Address()
		if err != nil {
			hlog.Fatal("get local ip failed: %s", err.Error())
		}
	}
}

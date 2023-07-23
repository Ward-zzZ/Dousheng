package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Key  string `mapstructure:"key" json:"key"`
}

type MsgEncryptConfig struct {
	Key string `mapstructure:"key" json:"key"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	MysqlInfo    MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	RedisInfo    RedisConfig    `mapstructure:"redis" json:"redis"`
	ConsulInfo   ConsulConfig   `mapstructure:"consul" json:"consul"`
	MsgEncrypt   MsgEncryptConfig `mapstructure:"msg_encrypt" json:"msg_encrypt"`
	OtelInfo     OtelConfig     `mapstructure:"otel" json:"otel"`
}

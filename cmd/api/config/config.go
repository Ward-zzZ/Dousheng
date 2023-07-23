package config

type JwtConfig struct {
	SecretKey string `mapstructure:"secret_key" json:"secret_key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Key  string `mapstructure:"key" json:"key"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type ServerConfig struct {
	Name            string       `mapstructure:"name" json:"name"`
	Host            string       `mapstructure:"host" json:"host"`
	Port            int          `mapstructure:"port" json:"port"`
	ProxyURL        string       `mapstructure:"proxy" json:"proxy"`
	JwtInfo         JwtConfig    `mapstructure:"jwt" json:"jwt"`
	OtelInfo        OtelConfig   `mapstructure:"otel" json:"otel"`
	CommentSrvInfo  RPCSrvConfig `mapstructure:"comment_srv" json:"comment_srv"`
	FavoriteSrvInfo RPCSrvConfig `mapstructure:"favorite_srv" json:"favorite_srv"`
	RelationSrvInfo RPCSrvConfig `mapstructure:"relation_srv" json:"relation_srv"`
	UserSrvInfo     RPCSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	VideoSrvInfo    RPCSrvConfig `mapstructure:"video_srv" json:"video_srv"`
	MessageSrvInfo  RPCSrvConfig `mapstructure:"message_srv" json:"message_srv"`
}

type RPCSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

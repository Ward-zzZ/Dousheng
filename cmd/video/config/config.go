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

type RabbitMqConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Key  string `mapstructure:"key" json:"key"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type MinioConfig struct {
	MinioURL				string `mapstructure:"minio_url" json:"minio_url"`
	MinioPort				string `mapstructure:"minio_port" json:"minio_port"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key" json:"secret_access_key"`
	VideoBucket     string `mapstructure:"video_bucket" json:"video_bucket"`
	CoverBucket     string `mapstructure:"cover_bucket" json:"cover_bucket"`
}

type ServerConfig struct {
	Name            string            `mapstructure:"name" json:"name"`
	Host            string            `mapstructure:"host" json:"host"`
	MysqlInfo       MysqlConfig       `mapstructure:"mysql" json:"mysql"`
	RedisInfo       RedisConfig       `mapstructure:"redis" json:"redis"`
	RabbitMqInfo    RabbitMqConfig    `mapstructure:"rabbitmq" json:"rabbitmq"`
	OtelInfo        OtelConfig        `mapstructure:"otel" json:"otel"`
	MinioInfo       MinioConfig       `mapstructure:"minio" json:"minio"`
	UserSrvInfo     UserSrvConfig     `mapstructure:"user_srv" json:"user_srv"`
	RelationSrvInfo RelationSrvConfig `mapstructure:"relation_srv" json:"relation_srv"`
	CommentSrvInfo  CommentSrvConfig  `mapstructure:"comment_srv" json:"comment_srv"`
	FavoriteSrvInfo FavoriteSrvConfig `mapstructure:"favorite_srv" json:"favorite_srv"`
}

type CommentSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type FavoriteSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type UserSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type RelationSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

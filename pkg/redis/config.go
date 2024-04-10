package redis

type Config struct {
	MasterName string   `mapstructure:"masterName" json:"masterName" default:""`
	Addresses  []string `mapstructure:"addresses" json:"addresses" default:""`
	DBNumber   int      `mapstructure:"dbNumber" json:"dbNumber" default:"0"`
	Password   string   `mapstructure:"password" json:"-" default:""`
}

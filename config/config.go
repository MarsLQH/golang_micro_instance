package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type Config struct {
	Consul struct {
		Addr string
	}
	Mysql struct {
		User string
		Pwd  string
		Host string
		Port string
		DB   string
	}
	Redis struct {
		Addr string
		Pwd  string
		DB   int
	}
}

var def Config

func C() Config {
	return def
}

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("")
	viper.SetEnvPrefix("")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Error().Err(err).Msg("配置文件未找到")
			panic(err)
		} else {
			log.Error().Err(err).Msg("找到了配置文件,发生了其它错误")
			panic(err)
		}
	}

	err = viper.Unmarshal(&def)
	if err != nil {
		panic(err)
	}
	return
	readRemoteConfig()

	log.Print(def)
}

func readRemoteConfig() {
	// 从环境变量中读取consul地址
	viper.BindEnv("consul_addr")
	if s := viper.GetString("consul_addr"); s != "" {
		def.Consul.Addr = s
	}
	if def.Consul.Addr == "" {
		return
	}
	log.Debug().Msgf("consul addr is %s", def.Consul.Addr)

	// 从consul中读取配置
	viper.AddRemoteProvider("consul", def.Consul.Addr, "")
	viper.SetConfigType("json")
	err := viper.ReadRemoteConfig()
	if err != nil {
		log.Warn().Err(err).Msg("not found the key  at consul key/value")
		return
	}

	readMysqlFromRemoteConfig()
	readRedisFromRemoteConfig()
}

func readMysqlFromRemoteConfig() {
	var v string
	if v = viper.GetString("mysql_db"); v != "" {
		def.Mysql.DB = v
	}
	if v = viper.GetString("mysql_user"); v != "" {
		def.Mysql.User = v
	}
	if v = viper.GetString("mysql_pwd"); v != "" {
		def.Mysql.Pwd = v
	}
	if v = viper.GetString("mysql_host"); v != "" {
		def.Mysql.Host = v
	}
	if v = viper.GetString("mysql_port"); v != "" {
		def.Mysql.Port = v
	}
}

func readRedisFromRemoteConfig() {
	var v string
	if v = viper.GetString("redis_addr"); v != "" {
		def.Redis.Addr = v
	}
	if v = viper.GetString("redis_pwd"); v != "" {
		def.Redis.Pwd = v
	}
}

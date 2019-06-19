package util

import (
	"encoding/json"
	"io/ioutil"
)

type Configs map[string]interface{}
type Config struct {
	section string
}

var configFile []byte
var RedisConfigs RedisConfig
var HostsConfig = make(map[string]map[string]interface{})

func init() {
	configFile, err = ioutil.ReadFile("config/Cmd.json")
	if err != nil {
		panic(err)
	}

	LoadConfig()
	LoadRedis(RedisConfigs.host, RedisConfigs.pass, RedisConfigs.port)
}

type RedisConfig struct {
	host string
	port float64
	pass string
}

func LoadConfig() bool {
	configs := Configs{}
	err = json.Unmarshal(configFile, &configs)
	if err != nil {
		panic(err)
	}

	// redis configs
	redisConfigMap := configs["redis"].(map[string]interface{})
	RedisConfigs.host = redisConfigMap["host"].(string)
	RedisConfigs.pass = redisConfigMap["pass"].(string)
	RedisConfigs.port = redisConfigMap["port"].(float64)

	// remote configs
	hostsMap := configs["hosts"].(map[string]interface{})
	for k, v := range hostsMap {
		HostsConfig[k] = v.(map[string]interface{})
	}

	return true
}

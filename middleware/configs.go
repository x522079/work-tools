package middleware

import (
	"cmds/resolve"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type Configs map[string]interface{}
type Config struct {
	section string
}

var ConfigFile []byte
var RedisConfigs RedisConfig
var HostsConfig = make(map[string]map[string]interface{})
var RedisUtil resolve.Redis
var RootDir string

type RedisConfig struct {
	host string
	port float64
	pass string
}

func init() {
	RootDir = os.Getenv("GO_CMD_DIR")
	config, err := ioutil.ReadFile(strings.TrimRight(RootDir, "/") + "/config/cmd.json")
	if err != nil {
		panic(err)
	}
	ConfigFile = config
}

func (c *Configs) LoadConfig() bool {
	configs := Configs{}
	err := json.Unmarshal(ConfigFile, &configs)
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

	RedisUtil.LoadRedis(RedisConfigs.host, RedisConfigs.pass, RedisConfigs.port)

	return true
}

func (c *Configs) SetRedisUtil(ru resolve.Redis) {
	RedisUtil = ru
}

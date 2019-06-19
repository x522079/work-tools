package util

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

var Conn redis.Conn
var err error

func LoadRedis(host, pass string, port float64) bool {
	Conn, err = redis.Dial("tcp", host+":"+strconv.Itoa(int(port)), redis.DialPassword(pass))
	if err != nil {
		panic(err)
	}

	return true
}

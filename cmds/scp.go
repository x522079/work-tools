package cmds

import (
	"cmds/util"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"strconv"
	"strings"
)

type Scp struct {
	SubCmd string
}

// 添加路径到内存
func (s *Scp) Add() {
	if len(Dispatcher.Cmd.Path) == 0 {
		panic("无效路径")
	}

	fmt.Println(getKey())
	res, err := redis.Int(util.Conn.Do("sadd", getKey(), Dispatcher.Cmd.Path))
	if err != nil {
		panic(err)
	}

	if res > 0 {
		fmt.Println("success")
	} else {
		fmt.Println("fail")
	}
}

// 列出所有路径
func (s *Scp) List() []string {
	key := getKey()
	res, err := redis.Values(util.Conn.Do("smembers", key))
	if err != nil {
		panic(err)
	}

	var list []string
	for _, v := range res {
		one := string(v.([]byte))
		list = append(list, one)
		fmt.Println(one)
	}

	return list
}

// 清空缓存
func (s *Scp) Clear() {
	res, err := redis.Int(util.Conn.Do("del", getKey()))
	if err != nil {
		panic(err)
	}
	if res == 0 {
		fmt.Println("缓存为空")
	}

	fmt.Println(res)
}

func (s *Scp) Execute() bool {
	paths := s.List()
	if len(paths) == 0 {
		fmt.Println("缓存为空，无需上传")
	}

	localClient := util.Client{}
	localClient.SetAlias(Dispatcher.Cmd.RemoteAlias)
	client := localClient.NewClient()
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	for _, v := range paths {
		info, err := os.Stat(v)
		if err != nil {
			panic(err)
		}

		fullRemotePath := localClient.GetFullHostPath() + ":" + Dispatcher.Cmd.RemotePath
		cmd := ""
		if info.IsDir() {
			cmd = "scp -r " + v + " " + fullRemotePath
		} else {
			cmd = "scp " + v + " " + fullRemotePath
		}

		if cmd == "" {
			continue
		}
		res, err := session.CombinedOutput(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(res)
	}
	return true
}

func (s *Scp) SubDispatch() bool {
	switch strings.ToLower(Dispatcher.Cmd.SubCmd) {
	case "add":
		s.Add()
		break
	case "clear":
		s.Clear()
		break
	case "list":
		s.List()
		break
	case "submit":
		s.Execute()
		break
	default:
		panic("不支持的子命令")
	}

	return true
}

func getKey() string {
	return "scp_" + strconv.Itoa(Dispatcher.Uid)
}

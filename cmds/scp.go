package cmds

import (
	"cmds/middleware"
	"cmds/util"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"os"
	"path/filepath"
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
	res, err := redis.Int(middleware.Conn.Do("sadd", getKey(), Dispatcher.Cmd.Path))
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
	res, err := redis.Values(middleware.Conn.Do("smembers", key))
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
	res, err := redis.Int(middleware.Conn.Do("del", getKey()))
	if err != nil {
		panic(err)
	}
	if res == 0 {
		fmt.Println("缓存为空")
	}

	fmt.Println(res)
}

/*func (s *Scp) Execute() bool {
	paths := s.List()
	if len(paths) == 0 {
		fmt.Println("缓存为空，无需上传")
	}

	localClient := util.Client{}
	localClient.SetAlias(Dispatcher.Cmd.RemoteAlias)
	localClient.Port = int(HostConfig["port"].(float64))
	localClient.User = HostConfig["user"].(string)
	localClient.Pass = HostConfig["pass"].(string)
	localClient.KeyPath = HostConfig["keyPath"].(string)
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

		fullRemotePath := localClient.GetFullHostPath(false) + ":" + Dispatcher.Cmd.RemotePath
		fmt.Println(fullRemotePath)
		cmd := ""
		if info.IsDir() {
			cmd = "scp -r " + v + " root@" + fullRemotePath
		} else {
			cmd = "scp " + v + " root@" + fullRemotePath
		}

		if cmd == "" {
			continue
		}

		res, err := session.CombinedOutput(cmd)
		if err != nil {
			fmt.Println(err)
			fmt.Println(string(res))
			continue
		}
		fmt.Println(string(res))
	}
	return true
}*/
func (s *Scp) Execute() bool {
	paths := s.List()
	if len(paths) == 0 {
		fmt.Println("缓存为空，无需上传")
	}

	localClient := util.Client{}
	localClient.SetAlias(Dispatcher.Cmd.RemoteAlias)
	localClient.Port = int(HostConfig["port"].(float64))
	localClient.User = HostConfig["user"].(string)
	localClient.Pass = HostConfig["pass"].(string)
	localClient.KeyPath = HostConfig["keyPath"].(string)
	sftpClient := localClient.NewSftpClient()
	for _, v := range paths {
		fmt.Println(v)
		file, err := os.Open(v)
		util.CheckError(err)
		defer file.Close()

		fullRemotePath := Dispatcher.Cmd.RemotePath + filepath.Base(file.Name())
		fmt.Println(fullRemotePath)
		remote, err := sftpClient.Create(fullRemotePath)
		util.CheckError(err)
		defer remote.Close()

		data, err := ioutil.ReadAll(file)
		util.CheckError(err)
		_, err = remote.Write(data)
		util.CheckError(err)
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

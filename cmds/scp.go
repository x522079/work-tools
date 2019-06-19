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
	"sync"
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
	g := sync.WaitGroup{}
	for _, v := range paths {
		g.Add(1)
		go func(v string) {
			defer g.Done()
			one := make([]string, 10)
			one = append(one, v)
			info, err := os.Stat(v)
			util.CheckError(err)

			if info.IsDir() {
				_ = filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}
					one = append(one, path)
					return nil
				})
			}

			gSub := sync.WaitGroup{}
			for _, vv := range one {
				gSub.Add(1)
				go func(vv string) {
					defer gSub.Done()

					file, err := os.Open(vv)
					if err != nil {
						return
					}
					defer file.Close()

					fullRemotePath := strings.TrimRight(Dispatcher.Cmd.RemotePath, "/") + "/" + filepath.Base(file.Name())
					remote, err := sftpClient.Create(fullRemotePath)
					util.CheckError(err)
					defer remote.Close()

					data, err := ioutil.ReadAll(file)
					util.CheckError(err)
					_, err = remote.Write(data)
					util.CheckError(err)
					fmt.Println(v, " >>> ", HostConfig["host"].(string)+"::"+fullRemotePath, " [ OK ]")
				}(vv)
			}

			gSub.Wait()
		}(v)
	}

	g.Wait()

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

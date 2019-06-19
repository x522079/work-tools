package main

import (
	"cmds/cmds"
	"flag"
	"fmt"
	"strings"
)

var commmand cmds.Commands

func init() {
	flag.StringVar(&commmand.Cmd, "cmd", "", "要执行的主命令")
	flag.StringVar(&commmand.SubCmd, "sub", "", "要执行的子命令")
	flag.StringVar(&commmand.Path, "path", "", "要添加的目录")
	flag.StringVar(&commmand.RemoteAlias, "remote", "", "远程服务器别名")
	var rpath = flag.String("rpath", "~", "远程服务器路径")
	flag.Parse()
	if strings.Contains(*rpath, "~") {
		*rpath = strings.Replace(*rpath, "~", "/root", 1)
	}
	commmand.RemotePath = *rpath
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			return
		}
	}()

	cmds.NewDispacher(commmand)
	cmds.Dispatcher.Dispach()
}

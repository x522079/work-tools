package main

import (
	"cmds/cmds"
	"flag"
	"fmt"
)

var commmand cmds.Commands

func init() {
	flag.StringVar(&commmand.Cmd, "cmd", "", "要执行的主命令")
	flag.StringVar(&commmand.SubCmd, "sub", "", "要执行的子命令")
	flag.StringVar(&commmand.Path, "path", "", "要添加的目录")
	flag.StringVar(&commmand.RemoteAlias, "remote", "", "远程服务器别名")
	flag.StringVar(&commmand.RemotePath, "rpath", "", "远程服务器路径")
}
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			return
		}
	}()

	flag.Parse()
	cmds.NewDispacher(commmand)
	cmds.Dispatcher.Dispach()
}

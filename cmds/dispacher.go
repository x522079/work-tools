package cmds

import (
	"cmds/middleware"
	_ "cmds/middleware"
	"math/rand"
)

type Commands struct {
	Cmd         string
	SubCmd      string
	Path        string
	RemotePath  string
	RemoteAlias string
}

type Dispacher struct {
	// 随机全局ID
	Uid int
	Cmd Commands
}

var (
	HostConfig map[string]interface{}
)

func NewDispacher(cmds Commands) {
	Dispatcher = &Dispacher{Uid: rand.Int(), Cmd: cmds}
}

func (d *Dispacher) Dispach() {
	HostConfig = middleware.HostsConfig[d.Cmd.RemoteAlias]

	switch d.Cmd.Cmd {
	case "scp":
		(&Scp{}).SubDispatch()
		break
	default:
		panic("暂不支持的命令")
	}
}
